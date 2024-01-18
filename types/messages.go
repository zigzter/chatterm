package types

type MessageKVMap struct {
	BadgeInfo        string `json:"badge-info"`
	Message          string `json:"message"`
	ID               string `json:"id"` // The UUID of the message itself
	Badges           string `json:"badges"`
	Color            string `json:"color"`
	DisplayName      string `json:"display-name"`
	Emotes           string `json:"emotes"`
	Login            string `json:"login"`
	Mod              string `json:"mod"`
	VIP              string `json:"vip"`
	MsgType          string `json:"msg-id"`
	CumulativeMonths string `json:"msg-param-cumulative-months"`   // Total months
	GiftMonths       string `json:"msg-param-months"`              // Total months if gifted
	ShareStreak      string `json:"msg-param-should-share-streak"` // Whehter the user wants their streak shared
	StreakMonths     string `json:"msg-param-streak-months"`
	SystemMsg        string `json:"system-msg"`  // Message shown to Twitch chat
	Timestamp        string `json:"tmi-sent-ts"` // Formatted timestamp
	ViewerCount      string `json:"msg-param-viewerCount"`
	ReceiverName     string `json:"msg-param-recipient-display-name"`
	GiftAmount       string `json:"msg-param-mass-gift-count"`
}

type Message interface {
	Implements()
}

type ChatMessage struct {
	Timestamp      string
	Color          string
	DisplayName    string
	IsFirstMessage bool
	IsMod          bool
	IsVIP          bool
	Message        string
	UserId         string
}

type SubMessage struct {
	DisplayName string
	Color       string
	Message     string
	Months      string
	Streak      string
	Timestamp   string
}

type SubGiftMessage struct {
	GiverName    string
	Color        string
	ReceiverName string
	Timestamp    string
}

type MysterySubGiftMessage struct {
	GiverName  string
	Color      string
	GiftAmount string
	Timestamp  string
}

type RaidMessage struct {
	DisplayName string
	Color       string
	ViewerCount string
	Timestamp   string
}

type AnnouncementMessage struct {
	DisplayName string
	Color       string
	Message     string
	Timestamp   string
}

type UserListMessage struct {
	Users []string
}

type UserStateMessage struct {
	Color         string
	IsMod         bool
	IsBroadcaster bool
}

type RoomStateMessage struct {
	ChannelID     string
	EmoteOnly     bool
	FollowersOnly bool
	SubsOnly      bool
	UniqueOnly    bool // Whether subsequent messages must be unique
	Slow          int  // Amount of seconds users must wait to send new messages
}

func (cm ChatMessage) Implements() {}

func (sm SubMessage) Implements() {}

func (sm SubGiftMessage) Implements() {}

func (mg MysterySubGiftMessage) Implements() {}

func (rm RaidMessage) Implements() {}

func (ul UserListMessage) Implements() {}

func (us UserStateMessage) Implements() {}

func (us RoomStateMessage) Implements() {}

func (am AnnouncementMessage) Implements() {}

type ParsedIRCMessage struct {
	Msg Message
}
