package utils

import (
	"strings"

	"github.com/zigzter/chatterm/types"
)

func RoomStateParser(input string) types.RoomStateMessage {
	roomState := types.RoomStateMessage{}
	parts := strings.SplitN(input, ":", 2)
	metadata := parts[0]
	keyValPairs := strings.Split(metadata, ";")
	for _, kvPair := range keyValPairs {
		kv := strings.Split(kvPair, "=")
		if len(kv) == 2 {
			key := kv[0]
			value := strings.TrimSpace(kv[1])
			switch key {
			case "room-id":
				roomState.ChannelID = value
			case "emote-only":
				roomState.EmoteOnly = value == "1"
			case "followers-only":
				roomState.FollowersOnly = value == "1"
			case "r9k":
				roomState.UniqueOnly = value == "1"
			case "subs-only":
				roomState.SubsOnly = value == "1"
			case "slow":
			}
		}
	}
	return roomState
}
