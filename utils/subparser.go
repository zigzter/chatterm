package utils

import (
	"strings"

	"github.com/zigzter/chatterm/types"
)

func SubParser(input string) types.SubMessage {
	parts := strings.SplitN(input, " :", 2)
	metadata := parts[0]
	message := ""
	messageSplit := strings.Split(parts[1], " :")
	if len(messageSplit) > 1 {
		message = messageSplit[1]
	}
	subMessage := types.SubMessage{Message: message}
	keyValPairs := strings.Split(metadata, ";")
	// TODO: Fix crash on no-message subs and gifts
	for _, kvPair := range keyValPairs {
		kv := strings.Split(kvPair, "=")
		if len(kv) == 2 {
			key := kv[0]
			value := kv[1]
			switch key {
			case "display-name":
				subMessage.DisplayName = value
			case "msg-param-cumulative-months":
				subMessage.Months = value
			case "msg-param-streak-months":
				subMessage.Streak = value
			}
		}
	}
	return subMessage
}
