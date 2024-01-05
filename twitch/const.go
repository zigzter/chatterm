package twitch

import "github.com/zigzter/chatterm/types"

const ClientId = "x6pl99d1tq9mqys6y2bmr59ahw9nik"

type CommandDetails struct {
	Endpoint string
	Method   string
}

var RequestMap = map[types.TwitchCommand]CommandDetails{
	types.Ban: {
		Endpoint: "/moderation/bans",
		Method:   "POST",
	},
	types.Clear: {
		Endpoint: "/moderation/chat",
		Method:   "DELETE",
	},
}
