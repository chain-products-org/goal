package twitter

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gophero/goal/redisx"
	"github.com/gophero/goal/stringx"

	"github.com/pkg/errors"
)

var EmptyAccessToken = AccessToken{}

type StateMap interface {
	Put(key string, state string)

	Get(key string) string

	Del(key string) string
}

type LocalStateMap struct {
	m map[string]*string
}

func NewLocalStateMap() StateMap {
	return &LocalStateMap{m: make(map[string]*string)}
}

func (l *LocalStateMap) Put(key string, state string) {
	l.m[key] = &state
}

func (l *LocalStateMap) Get(key string) string {
	r := l.m[key]
	if r == nil {
		return ""
	}
	return *r
}

func (l *LocalStateMap) Del(key string) string {
	r := l.m[key]
	delete(l.m, key)
	if r == nil {
		return ""
	}
	return *r
}

type RedisStateMap struct {
	key    string
	client redisx.Client
}

func NewRedisStateMap(c redisx.Client) StateMap {
	return &RedisStateMap{client: c, key: fmt.Sprintf("tw:st:mp:%d", time.Now().Unix())}
}

func (r *RedisStateMap) Put(key string, state string) {
	r.client.HSetNX(context.TODO(), r.key, key, state)
}

func (r *RedisStateMap) Get(key string) string {
	sc := r.client.HGet(context.TODO(), r.key, key)
	return sc.Val()
}

func (r *RedisStateMap) Del(key string) string {
	r.client.HDel(context.TODO(), r.key, key)
	return ""
}

type OAuth2AuthApi struct {
	sm StateMap
}

func NewOAuth2AuthApi(sm StateMap) *OAuth2AuthApi {
	if sm == nil {
		sm = NewLocalStateMap()
	}
	return &OAuth2AuthApi{sm: sm}
}

func (o *OAuth2AuthApi) AuthorizeUrl(clientId, redirectUri string, scope ...Scope) string {
	scopes := formatScopes(scope...)
	// code_challenge := stringx.Randn(16)
	state := stringx.Randn(16)
	// o.sm.Put(formatKey(clientId, state), code_challenge)
	// return fmt.Sprintf(auth2AuthorizeUrlFormat, clientId, redirectUri, scopes, state, code_challenge)
	// use state as challenge code
	return fmtUrl(auth2AuthorizeUrlFormat, clientId, redirectUri, scopes, state, state)
}

func (o *OAuth2AuthApi) tokenUrl() string {
	return fmtUrl(oauth2ApiUrlFormat, "/oauth2/token")
}

func (o *OAuth2AuthApi) revokeTokenUrl() string {
	return fmtUrl(oauth2ApiUrlFormat, "/oauth2/revoke")
}

func (o *OAuth2AuthApi) encodeClient(clientId, clientSecret string) string {
	dest := clientId + ":" + clientSecret
	return base64.StdEncoding.EncodeToString([]byte(dest))
}

func formatKey(clientId, key string) string {
	return clientId + "_" + key
}

func (o *OAuth2AuthApi) RequestAccessToken(clientId, clientSecret, code, state, redirectUri string) (AccessToken, error) {
	// challengeCode := o.sm.Get(formatKey(clientId, state))
	// o.sm.Del(formatKey(clientId, state))
	// if challengeCode == "" { // TODO not graceful
	// logx.Default.Errorf("challengecode is empty: %s", formatKey(clientId, state))
	// return EmptyAccessToken, errors.Wrapf(ApiError, "invalid state")
	// }
	body := strings.NewReader(
		NewGetParam().
			Append("code", code).
			Append("grant_type", "authorization_code").
			Append("redirect_uri", redirectUri).
			// Append("code_verifier", challengeCode).
			Append("code_verifier", state).
			Param(),
	)
	req, err := http.NewRequest(http.MethodPost, o.tokenUrl(), body)
	if err != nil {
		return EmptyAccessToken, errors.Wrapf(ErrApi, "request error: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+o.encodeClient(clientId, clientSecret))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return EmptyAccessToken, errors.Wrapf(ErrApi, "request token error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return AccessToken{}, errors.Wrapf(ErrApi, "request token not ok: %v", resp.Status)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return EmptyAccessToken, err
	}
	var accessToken AccessToken
	if err := json.Unmarshal(bs, &accessToken); err != nil {
		return EmptyAccessToken, errors.Wrapf(ErrApi, "invalid response: %v", string(bs))
	}
	return accessToken, nil
}

func (o *OAuth2AuthApi) RefreshAccessToken(clientId, clientSecret, refreshToken string) (AccessToken, error) {
	builder := stringx.NewBuilder()
	builder.WriteString("refresh_token=").WriteString(refreshToken).WriteString("&")
	builder.WriteString("grant_type=refresh_token").WriteString("&")
	builder.WriteString("client_id=").WriteString(clientId)
	body := strings.NewReader(builder.String())
	req, err := http.NewRequest(http.MethodPost, o.tokenUrl(), body)
	if err != nil {
		return EmptyAccessToken, errors.Wrapf(ErrApi, "request error: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+o.encodeClient(clientId, clientSecret))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return EmptyAccessToken, errors.Wrapf(ErrApi, "request error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return AccessToken{}, errors.Wrapf(ErrApi, "request error: %v", resp.Status)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return EmptyAccessToken, err
	}
	var accessToken AccessToken
	if err := json.Unmarshal(bs, &accessToken); err != nil {
		return EmptyAccessToken, errors.Wrapf(ErrApi, "invalid response: %v", string(bs))
	}
	return accessToken, nil
}

func (o *OAuth2AuthApi) RevokeAccessToken(clientId, clientSecret string, token string) error {
	builder := stringx.NewBuilder()
	builder.WriteString("token=").WriteString(token).WriteString("&")
	builder.WriteString("token_type_hint=access_token")
	body := strings.NewReader(builder.String())
	req, err := http.NewRequest(http.MethodPost, o.revokeTokenUrl(), body)
	if err != nil {
		return errors.Wrapf(ErrApi, "request error: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+o.encodeClient(clientId, clientSecret))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(ErrApi, "request error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Wrapf(ErrApi, "request error: %v", resp.Status)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var revokeResut map[string]bool
	if err := json.Unmarshal(bs, &revokeResut); err != nil {
		return errors.Wrapf(ErrApi, "invalid response: %v", string(bs))
	}
	if !revokeResut["revoked"] {
		return fmt.Errorf("revoke failed")
	}
	return nil
}

type AccessToken struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

// AppOnlyToken to get app-only token through consumer keys(api key and secret).
// see: https://developer.x.com/en/docs/authentication/oauth-2-0/application-only
func (o *OAuth2AuthApi) AppOnlyToken(appId, appSecret string) (AccessToken, error) {
	// encClientId := base64.URLEncoding.EncodeToString([]byte(clientId))
	// encClientSecret := base64.URLEncoding.EncodeToString([]byte(clientSecret))

	builder := stringx.NewBuilder()
	builder.WriteString("grant_type=").WriteString("client_credentials")
	body := strings.NewReader(builder.String())

	url := "https://api.twitter.com/oauth2/token"
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return EmptyAccessToken, errors.Wrapf(ErrApi, "request error: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+o.encodeClient(appId, appSecret))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return EmptyAccessToken, errors.Wrapf(ErrApi, "request token error: %v", err)
	}
	bs, err := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return AccessToken{}, errors.Wrapf(ErrApi, "request token not ok: %v", resp.Status)
	}

	if err != nil {
		return EmptyAccessToken, err
	}
	var accessToken AccessToken
	if err := json.Unmarshal(bs, &accessToken); err != nil {
		return EmptyAccessToken, errors.Wrapf(ErrApi, "invalid response: %v", string(bs))
	}
	return accessToken, nil
}
