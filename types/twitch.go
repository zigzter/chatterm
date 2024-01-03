package types

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