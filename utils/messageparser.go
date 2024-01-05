package utils

import (
	"strconv"
	"strings"
	"time"

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
			case "mod":
				chatMessage.IsMod = value == "1"
			case "vip":
				chatMessage.IsVIP = value == "1"
			case "tmi-sent-ts":
				unixTime, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					chatMessage.Timestamp = "00:00"
					return chatMessage
				}
				timeObj := time.Unix(unixTime/1000, 0)
				// TODO: Custom timezone
				location, err := time.LoadLocation("America/Los_Angeles")
				if err != nil {
					chatMessage.Timestamp = "00:00"
					return chatMessage
				}
				localTime := timeObj.In(location)
				formattedTime := localTime.Format("15:04")
				chatMessage.Timestamp = formattedTime
			}
		}
	}
	return chatMessage
}
