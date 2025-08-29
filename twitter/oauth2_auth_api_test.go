package twitter_test

import (
	"fmt"
	"testing"

	"github.com/chain-products-org/goal/twitter"
)

var (
	scopes = "offline.access%20tweet.read%20users.read%20follows.read%20follows.write"
	url    = "https://twitter.com/i/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=state&code_challenge=%s&code_challenge_method=plain"
)

func TestRequestToken(t *testing.T) {
	state := ""
	redirectUri := "http://localhost:8080"
	ak, err := twitter.OAuth2Apis.Auth.RequestAccessToken(testEnv.setting.ClientId, testEnv.setting.ClientSecret, testEnv.code, state, redirectUri)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ak)
}

func TestRefreshToken(t *testing.T) {
	refreshToken := testEnv.refrshToken
	ss, err := twitter.OAuth2Apis.Auth.RefreshAccessToken(testEnv.setting.ClientId, testEnv.setting.ClientSecret, refreshToken)
	if err != nil {
		t.Error(err)
	}
	testEnv.accessToken = ss.AccessToken
	testEnv.refrshToken = ss.RefreshToken
}

func TestRevokeToken(t *testing.T) {
	token := testEnv.accessToken
	err := twitter.OAuth2Apis.Auth.RevokeAccessToken(testEnv.setting.ClientId, testEnv.setting.ClientSecret, token)
	if err != nil {
		t.Error(err)
	}
}

func TestAppOnlyToken(t *testing.T) {
	at, err := twitter.OAuth2Apis.Auth.AppOnlyToken(testEnv.setting.ApiKey, testEnv.setting.ApiSecret)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(at)
}
