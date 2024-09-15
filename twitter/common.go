package twitter

import (
	"fmt"
	"strings"

	"github.com/gophero/goal/collection/slicex"
)

var ErrApi = fmt.Errorf("twitter api error")

type Error struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

type Result[T any] struct {
	Data T    `json:"data"`
	Meta Meta `json:"meta,omitempty"`
}

type Meta struct {
	ResultCount   uint32 `json:"result_count"`
	PreviousToken string `json:"previous_token"`
	NextToken     string `json:"next_token"`
}

type FieldFilter struct {
	Expansions    []Expansion
	TwitterFields []TwitterField
	UserFields    []UserField
}

func NewFieldFilter() *FieldFilter {
	return &FieldFilter{}
}

func (ff *FieldFilter) AddExpansion(exps ...Expansion) *FieldFilter {
	ff.Expansions = append(ff.Expansions, exps...)
	return ff
}

func (ff *FieldFilter) AddTwitterField(tfs ...TwitterField) *FieldFilter {
	ff.TwitterFields = append(ff.TwitterFields, tfs...)
	return ff
}

func (ff *FieldFilter) AddUserField(ufs ...UserField) *FieldFilter {
	ff.UserFields = append(ff.UserFields, ufs...)
	return ff
}

type (
	Expansion    string
	TwitterField string
	UserField    string
)

const (
	ExpansionPinnedTweetId Expansion = "pinned_tweet_id"
)

const (
	TwitterFieldAttachments        TwitterField = "attachments"
	TwitterFieldAuthorId           TwitterField = "author_id"
	TwitterFieldContextAnnotations TwitterField = "context_annotations"
	TwitterFieldConversationId     TwitterField = "conversation_id"
	TwitterFieldCreatedAt          TwitterField = "created_at"
	TwitterFieldEditControls       TwitterField = "edit_controls"
	TwitterFieldEntities           TwitterField = "entities"
	TwitterFieldGeo                TwitterField = "geo"
	TwitterFieldId                 TwitterField = "id"
	TwitterFieldInReplyToUserId    TwitterField = "in_reply_to_user_id"
	TwitterFieldLang               TwitterField = "lang"
	TwitterFieldNonPublicMetrics   TwitterField = "non_public_metrics"
	TwitterFieldPublicMetrics      TwitterField = "public_metrics"
	TwitterFieldOrganicMetrics     TwitterField = "organic_metrics"
	TwitterFieldPromotedMetrics    TwitterField = "promoted_metrics"
	TwitterFieldPossiblySensitive  TwitterField = "possibly_sensitive"
	TwitterFieldReferencedTweets   TwitterField = "referenced_tweets"
	TwitterFieldReplySettings      TwitterField = "reply_settings"
	TwitterFieldSource             TwitterField = "source "
	TwitterFieldText               TwitterField = "text"
	TwitterFieldWithheld           TwitterField = "withheld"
)

const (
	UserFieldCreatedAt         UserField = "created_at"
	UserFieldDescription       UserField = "description"
	UserFieldEntities          UserField = "entities"
	UserFieldId                UserField = "id"
	UserFieldLocation          UserField = "location"
	UserFieldMostRecentTweetId UserField = "most_recent_tweet_id"
	UserFieldName              UserField = "name"
	UserFieldPinnedTweetId     UserField = "pinned_tweet_id"
	UserFieldProfileImageUrl   UserField = "profile_image_url"
	UserFieldProtected         UserField = "protected"
	UserFieldPublicMetrics     UserField = "public_metrics"
	UserFieldUrl               UserField = "url"
	UserFieldUserName          UserField = "username"
	UserFieldVerified          UserField = "verified"
	UserFieldVerifiedType      UserField = "verified_type"
	UserFieldWithHeld          UserField = "withheld"
)

func formatExpansion(exps ...Expansion) string {
	rs := slicex.Eachv(exps, func(v Expansion) string {
		return string(v)
	})
	return strings.Join(rs, ",")
}

func formatTweetFields(tfs ...TwitterField) string {
	rs := slicex.Eachv(tfs, func(v TwitterField) string {
		return string(v)
	})
	return strings.Join(rs, ",")
}

func formatUserFields(ufs ...UserField) string {
	rs := slicex.Eachv(ufs, func(v UserField) string {
		return string(v)
	})
	return strings.Join(rs, ",")
}
