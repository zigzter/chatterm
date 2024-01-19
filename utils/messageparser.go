package utils

import (
	"strings"

	"github.com/zigzter/chatterm/types"
)

func MessageParser(input string) types.ChatMessage {
	parts := strings.SplitN(input, " :", 2)
	metadata := parts[0]
	message := strings.Split(parts[1], " :")[1]
	keyValPairs := strings.Split(metadata, ";")
	chatMessage := types.ChatMessage{Message: message}
	for _, kvPair := range keyValPairs {
		kv := strings.Split(kvPair, "=")
		if len(kv) == 2 {
			key := kv[0]
			value := kv[1]

			switch key {
			case "user-id":
				chatMessage.UserId = value
			case "color":
				chatMessage.Color = value
			case "display-name":
				chatMessage.DisplayName = value
			case "first-msg":
				chatMessage.IsFirstMessage = value == "1"
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
				chatMessage.ChannelUserType = channelUserType
			case "tmi-sent-ts":
				chatMessage.Timestamp = ParseTimestamp(value)
			}
		}
	}
	return chatMessage
}
