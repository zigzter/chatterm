package types

type AuthResultMsg struct {
	Success bool
	Error   string
	Token   string
}
