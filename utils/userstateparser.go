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
			case "mod":
				userState.IsMod = value == "1"
			case "badges":
				badges := strings.Split(value, ",")
				isBroadcaster := false
				for _, badge := range badges {
					if strings.HasPrefix(badge, "broadcaster") && strings.HasSuffix(value, "/1") {
						isBroadcaster = true
					}
				}
				userState.IsBroadcaster = isBroadcaster
			}
		}
	}
	return userState
}
