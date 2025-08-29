package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/chain-products-org/goal/twitter"
)

var temp *template.Template

func init() {
	// it's relative to the project's root director
	tmpl, err := template.New("demo.html").ParseFiles("twitter/demo/demo.html")
	if err != nil {
		panic(err)
	}
	temp = tmpl
}

func main() {
	clientId, clientSecret, redirectUrl := "bW55cThSSFlGR3ktMWxxWUIxMzY6MTpjaQ", "", "http://localhost:8080/demo"
	clientSecret = os.Getenv("gophero_x_secret") // read client recret from env
	mux := http.DefaultServeMux
	mux.HandleFunc("/demo", func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
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
				"title":        "test html",
				"authUrl":      "",
				"redirect_uri": redirectUrl,
				"err":          err,
				"at":           at.AccessToken,
				"rt":           at.RefreshToken,
				"user":         userInfo,
			}
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			temp.Execute(w, data)
		} else {
			data := map[string]any{
				"title":        "test html",
				"authUrl":      twitter.OAuth2Apis.Auth.AuthorizeUrl(clientId, redirectUrl, twitter.TweetRead, twitter.TweetWrite, twitter.OfflineAccess, twitter.UsersRead, twitter.FollowsRead, twitter.FollowsWrite),
				"err":          err,
				"redirect_uri": redirectUrl,
			}
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			temp.Execute(w, data)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
