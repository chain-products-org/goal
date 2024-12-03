package twitter

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type OAuth2TweetApi struct{}

func NewOAuth2TweetApi() *OAuth2TweetApi {
	return &OAuth2TweetApi{}
}

// PostTweet Creates a Tweet on behalf of an authenticated user.
// OAuth 2.0 scopes need: tweet.read,tweet.write,users.read
func (o *OAuth2TweetApi) PostTweet(accessToken string, param PostTweetParam) (*PostTweetResp, error) {
	var err error
	url := fmtUrl(oauth2ApiUrlFormat, "/tweets")
	body := bytes.NewReader(param.Json())
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return &PostTweetResp{}, errors.Wrapf(ErrApi, "request error: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &PostTweetResp{}, errors.Wrapf(ErrApi, "request error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return &PostTweetResp{}, errors.Wrapf(ErrApi, "request error: %v", resp.Status)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return &PostTweetResp{}, err
	}

	var r Result[PostTweetResp]
	if err = json.Unmarshal(bs, &r); err != nil {
		return nil, err
	}
	return &r.Data, nil
}

func (o *OAuth2TweetApi) RetweetBy(accessToken, tweetId string, ff *FieldFilter, options ...GetParamOption) ([]*UserInfo, Meta, error) {
	url := fmtUrl(oauth2ApiUrlFormat, "/tweets/"+tweetId+"/retweeted_by")
	params := NewGetParam().FilterFields(ff)
	for _, p := range options {
		p(params)
	}
	q := params.Param()
	req, err := http.NewRequest(http.MethodGet, url+"?"+q, nil)
	if err != nil {
		return []*UserInfo{}, Meta{}, errors.Wrapf(ErrApi, "request error: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []*UserInfo{}, Meta{}, errors.Wrapf(ErrApi, "request error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return []*UserInfo{}, Meta{}, errors.Wrapf(ErrApi, "request error: %v", resp.Status)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return []*UserInfo{}, Meta{}, err
	}
	var result Result[[]*UserInfo]
	if err := json.Unmarshal(bs, &result); err != nil {
		return []*UserInfo{}, Meta{}, errors.Wrapf(ErrApi, "invalid response: %v", string(bs))
	}
	return result.Data, result.Meta, nil
}

type Tweet struct{}
