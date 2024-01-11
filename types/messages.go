package types

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
	Message     string
	Months      string
	Streak      string
}

type UserListMessage struct {
	Users []string
}

func (cm ChatMessage) Implements() {}

func (sm SubMessage) Implements() {}

func (ul UserListMessage) Implements() {}

type ParsedIRCMessage struct {
	Msg Message
}
