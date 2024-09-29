package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gophero/goal/twitter"
)

func main() {
	clientId, clientSecret, redirectUrl := "", "", ""
	mux := http.DefaultServeMux
	mux.HandleFunc("/demo", func(w http.ResponseWriter, req *http.Request) {
		err := req.Form.Get("error")   // if there is an error
		state := req.Form.Get("state") // state is using to prevent csrf attacks：https://auth0.com/docs/protocols/state-parameters
		code := req.Form.Get("code")   // auth code
		if code != "" && err == "" {
			// get accessToken
			at, err := twitter.OAuth2Apis.Auth.RequestAccessToken(clientId, clientSecret, code, state, redirectUrl)
			fmt.Println("redisClient id:", clientId)
			fmt.Println("access token :", at.AccessToken)
			fmt.Println("refresh token:", at.RefreshToken)
			bs, _ := json.Marshal(at)
			fmt.Println(string(bs))
			userInfo, err := twitter.OAuth2Apis.User.Me(at.AccessToken, twitter.NewFieldFilter().AddUserField(twitter.UserFieldProfileImageUrl))
			fmt.Println(userInfo)

			data := map[string]any{
				"title":   "test html",
				"authUrl": "",
				"err":     err,
				"at":      at.AccessToken,
				"rt":      at.RefreshToken,
				"user":    userInfo,
			}
			temp := template.New("demo.html")
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			temp.Execute(w, data)
		} else {
			data := map[string]any{
				"title":   "test html",
				"authUrl": twitter.OAuth2Apis.Auth.AuthorizeUrl(clientId, redirectUrl, twitter.TweetRead, twitter.TweetWrite, twitter.OfflineAccess, twitter.UsersRead, twitter.FollowsRead, twitter.FollowsWrite),
				"err":     err,
			}
			temp := template.New("demo.html")
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			temp.Execute(w, data)
		}
		http.ServeFile(w, req, "demo.html")
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
