package twitter_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gophero/goal/twitter"
	"github.com/stretchr/testify/assert"
)

func TestMe(t *testing.T) {
	user, err := twitter.OAuth2Apis.User.Me(testEnv.accessToken, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(user)
	assert.True(t, user.Id != "" && user.Name != "" && user.Username != "" && user.ProfileImageUrl == "")
	ff := twitter.NewFieldFilter()
	ff.AddUserField(twitter.UserFieldId, twitter.UserFieldProfileImageUrl, twitter.UserFieldCreatedAt, twitter.UserFieldVerified,
		twitter.UserFieldVerifiedType, twitter.UserFieldWithHeld, twitter.UserFieldDescription, twitter.UserFieldLocation,
		twitter.UserFieldPublicMetrics)
	user, err = twitter.OAuth2Apis.User.Me(testEnv.accessToken, ff)
	if err != nil {
		t.Fatal(err)
	}
	bs, _ := json.Marshal(user)
	fmt.Println(string(bs))
	assert.True(t, user.CreatedAt.Year() > 2000 && user.ProfileImageUrl != "")
}

func TestFollowers(t *testing.T) {
	token, err := twitter.OAuth2Apis.Auth.AppOnlyToken(testEnv.setting.ApiKey, testEnv.setting.ApiSecret)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	id := testEnv.id
	users, meta, err := twitter.OAuth2Apis.User.Followers(token.AccessToken, id, nil, twitter.GetParamOptions.MaxResults(2))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("users:", users)
	fmt.Println("meta:", meta)
	ff := twitter.NewFieldFilter()
	ff.AddUserField(twitter.UserFieldId, twitter.UserFieldProfileImageUrl, twitter.UserFieldCreatedAt, twitter.UserFieldVerified, twitter.UserFieldWithHeld, twitter.UserFieldDescription, twitter.UserFieldLocation)
	users, _, err = twitter.OAuth2Apis.User.Followers(testEnv.accessToken, id, ff, twitter.GetParamOptions.MaxResults(1000), twitter.GetParamOptions.PaginationToken(meta.NextToken))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(users)

	// {"client_id":"28659504","detail":"When authenticating requests to the Twitter API v2 endpoints, you must use keys and tokens from a Twitter developer App that is attached to a Project. You can create a project via the developer portal.","registration_url":"https://developer.twitter.com/en/docs/projects/overview","title":"Client Forbidden","required_enrollment":"Appropriate Level of API Access","reason":"client-not-enrolled","type":"https://api.twitter.com/2/problems/client-forbidden"}
}

func TestModel(t *testing.T) {
	js := `
{
  "data": [
    {
      "id": "6253282",
      "name": "Twitter API",
      "username": "TwitterAPI"
    },
    {
      "id": "2244994945",
      "name": "Twitter Dev",
      "username": "TwitterDev"
    },
    {
      "id": "783214",
      "name": "Twitter",
      "username": "Twitter"
    },
    {
      "id": "95731075",
      "name": "Twitter Safety",
      "username": "TwitterSafety"
    },
    {
      "id": "3260518932",
      "name": "Twitter Moments",
      "username": "TwitterMoments"
    },
    {
      "id": "373471064",
      "name": "Twitter Music",
      "username": "TwitterMusic"
    },
    {
      "id": "791978718",
      "name": "Twitter Official Partner",
      "username": "OfficialPartner"
    },
    {
      "id": "17874544",
      "name": "Twitter Support",
      "username": "TwitterSupport"
    },
    {
      "id": "234489024",
      "name": "Twitter Comms",
      "username": "TwitterComms"
    },
    {
      "id": "1526228120",
      "name": "Twitter Data",
      "username": "TwitterData"
    }
  ],
  "meta": {
    "result_count": 10,
    "next_token": "DFEDBNRFT3MHCZZZ"
  }
}
`
	var r twitter.Result[[]*twitter.UserInfo]
	if err := json.Unmarshal([]byte(js), &r); err != nil {
		t.Errorf("test failed: %v", err)
	} else {
		users := r.Data
		assert.True(t, r.Meta.ResultCount == 10)
		assert.True(t, uint32(len(users)) == r.Meta.ResultCount)
		assert.True(t, r.Meta.NextToken == "DFEDBNRFT3MHCZZZ")
	}
}
