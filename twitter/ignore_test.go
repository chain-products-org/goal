package twitter_test

import (
	"fmt"
	"os"

	"github.com/chain-products-org/goal/errorx"
	"github.com/chain-products-org/goal/twitter"
)

const (
	accessToken  = "x_accesstoken"
	refreshToken = "x_refreshtoken"
	bearerToken  = "x_bearertoken"
)

var testEnv testenv

type testenv struct {
	setting     twitter.Setting
	code        string
	accessToken string
	refrshToken string
	bearerToken string
	id          string
}

func init() {
	testEnv.setting = twitter.Setting{
		ClientId:     os.Getenv("x_clientid"),
		ClientSecret: os.Getenv("x_clientsecret"),
	}
	testEnv.code = ""
	testEnv.accessToken = os.Getenv("x_accesstoken")
	testEnv.refrshToken = os.Getenv("x_refreshtoken")
	testEnv.bearerToken = os.Getenv("x_bearertoken")
	testEnv.id = os.Getenv("x_id")
	testEnv.setting.ApiKey = os.Getenv("x_appid")
	testEnv.setting.ApiSecret = os.Getenv("x_app_secret")
	fmt.Println("testenv:", testEnv)
}

func resetToken(at, rt, bt string) {
	err := os.Setenv(accessToken, at)
	errorx.Throw(err)
	err = os.Setenv(refreshToken, rt)
	errorx.Throw(err)
	err = os.Setenv(bearerToken, bt)
	errorx.Throw(err)
}
