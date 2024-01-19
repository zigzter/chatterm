package utils

import (
	"strings"

	"github.com/zigzter/chatterm/types"
)

func RoomStateParser(input string) *types.RoomStateMessage {
	roomState := types.RoomStateMessage{}
	parts := strings.SplitN(input, ":", 2)
	metadata := parts[0]
	keyValPairs := strings.Split(metadata, ";")
	for _, kvPair := range keyValPairs {
		kv := strings.Split(kvPair, "=")
		if len(kv) == 2 {
			key := strings.TrimPrefix(kv[0], "@")
			value := strings.TrimSpace(kv[1])
			switch key {
			case "room-id":
				roomState.ChannelID = &value
			case "emote-only":
				enabled := value == "1"
				roomState.EmoteOnly = &enabled
			case "followers-only":
				enabled := value != "-1"
				roomState.FollowersOnly = &enabled
			case "r9k":
				enabled := value == "1"
				roomState.UniqueOnly = &enabled
			case "subs-only":
				enabled := value == "1"
				roomState.SubOnly = &enabled
			case "slow":
				roomState.Slow = &value
			}
		}
	}
	return &roomState
}
