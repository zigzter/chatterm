package types

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
