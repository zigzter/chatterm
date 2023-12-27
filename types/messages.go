package types

type Message interface {
	GetMainText() string
}

type ChatMessage struct {
	Timestamp      string
	Color          string
	DisplayName    string
	IsFirstMessage bool
	IsMod          bool
	IsVIP          bool
	Message        string
}

type SubMessage struct {
	DisplayName string
	Message     string
	Months      string
	Streak      string
}

func (cm ChatMessage) GetMainText() string {
	return cm.Message
}

func (sm SubMessage) GetMainText() string {
	return sm.Message
}

type ChatMessageWrap struct {
	ChatMsg Message
}
