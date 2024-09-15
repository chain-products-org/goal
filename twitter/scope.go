package twitter

import "strings"

type Scope string

const (
	TweetRead          Scope = "tweet.read"           // All the Tweets you can view, including Tweets from protected accounts.
	TweetWrite         Scope = "tweet.write"          // Tweet and Retweet for you.
	TweetModerateWrite Scope = "tweet.moderate.write" // Hide and unhide replies to your Tweets.
	UsersRead          Scope = "users.read"           // Any account you can view, including protected accounts.
	FollowsRead        Scope = "follows.read"         // People who follow you and people who you follow.
	FollowsWrite       Scope = "follows.write"        // Follow and unfollow people for you.
	OfflineAccess      Scope = "offline.access"       // Stay connected to your account until you revoke access.
	SpaceRead          Scope = "space.read"           // All the Spaces you can view.
	MuteRead           Scope = "mute.read"            // Accounts you’ve muted.
	MuteWrite          Scope = "mute.write"           // Mute and unmute accounts for you.
	LikeRead           Scope = "like.read"            // Tweets you’ve liked and likes you can view.
	LikeWrite          Scope = "like.write"           // Like and un-like Tweets for you.
	ListRead           Scope = "list.read"            // Lists, list members, and list followers of lists you’ve created or are a member of, including private lists.
	ListWrite          Scope = "list.write"           // Create and manage Lists for you.
	BlockRead          Scope = "block.read"           // Accounts you’ve blocked.
	BlockWrite         Scope = "block.write"          // Block and unblock accounts for you.
	BookmarkRead       Scope = "bookmark.read"        // Get Bookmarked Tweets from an authenticated user.
	BookmarkWrite      Scope = "bookmark.write"       // Bookmark and remove Bookmarks from Tweets
)

func formatScopes(scopes ...Scope) string {
	var sb strings.Builder
	for i, s := range scopes {
		sb.WriteString(string(s))
		if i < len(scopes)-1 {
			sb.WriteString("%20")
		}
	}
	return sb.String()
}
