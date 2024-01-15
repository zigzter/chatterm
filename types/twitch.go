package types

import (
	"fmt"
	"time"
)

type (
	ServerStartMsg     struct{}
	ServerStartedMsg   struct{}
	AuthOpenMsg        struct{}
	AuthOpenedMsg      struct{}
	TokenReceiveMsg    struct{}
	TokenReceivedMsg   struct{}
	ProcessCompleteMsg struct{}
)

type TwitchCommand string

const (
	Ban          TwitchCommand = "ban"
	Info         TwitchCommand = "info"
	Unban        TwitchCommand = "unban"
	Clear        TwitchCommand = "clear"
	Delete       TwitchCommand = "delete"
	Slow         TwitchCommand = "slow"
	SubOnly      TwitchCommand = "subonly"
	User         TwitchCommand = "user"
	Subscription TwitchCommand = "subscription"
	Followers    TwitchCommand = "followers"
	LiveChannels TwitchCommand = "livechannels"
)

type TwitchAPIError struct {
	ServerError string `json:"error"`
	Status      int    `json:"status"`
	Message     string `json:"message"`
}

func (e TwitchAPIError) Error() string {
	return fmt.Sprintf("Code %d, message: %s", e.Status, e.Message)
}

// UserInfo merges the data from several API calls together
type UserInfo struct {
	Details      UserData
	Following    FollowersData
	Subscription SubscriptionData
}

type UserData struct {
	ID              string    `json:"id"`
	Login           string    `json:"login"`
	DisplayName     string    `json:"display_name"`
	Type            string    `json:"type"`
	BroadcasterType string    `json:"broadcaster_type"`
	Description     string    `json:"description"`
	ProfileImageURL string    `json:"profile_image_url"`
	OfflineImageURL string    `json:"offline_image_url"`
	ViewCount       int       `json:"view_count"`
	CreatedAt       time.Time `json:"created_at"`
}

type UserResp struct {
	Data []UserData `json:"data"`
}

type SubscriptionData struct {
	BroadcasterID    string `json:"broadcaster_id"`
	BroadcasterName  string `json:"broadcaster_name"`
	BroadcasterLogin string `json:"broadcaster_login"`
	IsGift           bool   `json:"is_gift"`
	GifterName       string `json:"gifter_name"` // Only exists if IsGift
	Tier             string `json:"tier"`        // 1000, 2000, 3000
}

type SubscriptionResp struct {
	Data []SubscriptionData `json:"data"`
}

type FollowersData struct {
	ID          string `json:"user_id"`
	Displayname string `json:"user_login"`
	FollowedAt  string `json:"followed_at"`
}

type FollowersResp struct {
	Data []FollowersData `json:"data"`
}

type LiveChannelsData struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserLogin    string    `json:"user_login"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	GameName     string    `json:"game_name"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
	TagIds       []any     `json:"tag_ids"`
	Tags         []string  `json:"tags"`
}

type LiveChannelsResp struct {
	Data       []LiveChannelsData `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

type UserBanResp struct {
	Data []struct {
		ChannelID   string `json:"broadcaster_id"`
		ModeratorID string `json:"moderator_id"`
		UserID      string `json:"user_id"` // User ID of the ban target
		CreatedAt   string `json:"created_at"`
		EndTime     any    `json:"end_time"` // null if ban, string if timeout
	} `json:"data"`
}

type UpdateChatSettingsResp struct {
	Data []struct {
		ChannelID                     string `json:"broadcaster_id"`
		ModeratorID                   string `json:"moderator_id"`
		SlowMode                      bool   `json:"slow_mode"`
		SlowModeWaitTime              int    `json:"slow_mode_wait_time"`
		FollowerMode                  bool   `json:"follower_mode"`
		FollowerModeDuration          any    `json:"follower_mode_duration"`
		SubscriberMode                bool   `json:"subscriber_mode"`
		EmoteMode                     bool   `json:"emote_mode"`
		UniqueChatMode                bool   `json:"unique_chat_mode"`
		NonModeratorChatDelay         bool   `json:"non_moderator_chat_delay"`
		NonModeratorChatDelayDuration any    `json:"non_moderator_chat_delay_duration"`
	} `json:"data"`
}
