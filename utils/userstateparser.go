package utils

import (
	"strings"

	"github.com/zigzter/chatterm/types"
)

func UserStateParser(input string) types.UserStateMessage {
	userState := types.UserStateMessage{}
	parts := strings.SplitN(input, ":", 2)
	metadata := parts[0]
	keyValPairs := strings.Split(metadata, ";")
	for _, kvPair := range keyValPairs {
		kv := strings.Split(kvPair, "=")
		if len(kv) == 2 {
			key := kv[0]
			value := kv[1]
			switch key {
			case "color":
				userState.Color = value
			case "badges":
				// IRC message does have "mod" & "vip" fields, but it does not contain a broadcaster field,
				// so we're just doing it all within badges
				badges := strings.Split(value, ",")
				channelUserType := "normal"
				for _, badge := range badges {
					if strings.HasPrefix(badge, "broadcaster") && strings.HasSuffix(badge, "/1") {
						channelUserType = "broadcaster"
					}
					if strings.HasPrefix(badge, "moderator") && strings.HasSuffix(badge, "/1") {
						channelUserType = "moderator"
					}
					if strings.HasPrefix(badge, "vip") && strings.HasSuffix(badge, "/1") {
						channelUserType = "vip"
					}
				}
				userState.ChannelUserType = channelUserType
			}
		}
	}
	return userState
}
