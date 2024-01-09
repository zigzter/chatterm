package types

import "time"

type AuthResultMsg struct {
	Success bool
	Error   string
	Token   string
}

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
	Ban     TwitchCommand = "ban"
	Unban   TwitchCommand = "unban"
	Clear   TwitchCommand = "clear"
	Delete  TwitchCommand = "delete"
	Slow    TwitchCommand = "slow"
	SubOnly TwitchCommand = "subonly"
)

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

type FetchUserResp struct {
	Data []UserData `json:"data"`
}
