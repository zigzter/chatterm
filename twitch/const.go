package twitch

import (
	"github.com/zigzter/chatterm/types"
)

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
	types.Unban: {
		Endpoint: "/moderation/bans",
		Method:   "DELETE",
	},
	types.Slow: {
		Endpoint: "/chat/settings",
		Method:   "PATCH",
	},
	types.SubOnly: {
		Endpoint: "/chat/settings",
		Method:   "PATCH",
	},
	types.Delete: {
		Endpoint: "/moderation/chat",
		Method:   "DELETE",
	},
	types.Clear: {
		Endpoint: "/moderation/chat",
		Method:   "DELETE",
	},
	types.User: {
		Endpoint: "/users",
		Method:   "GET",
	},
	types.Subscription: {
		Endpoint: "/subscriptions/user",
		Method:   "GET",
	},
	types.Followers: {
		Endpoint: "/channels/followers",
		Method:   "GET",
	},
	types.LiveChannels: {
		Endpoint: "/streams/followed",
		Method:   "GET",
	},
}
