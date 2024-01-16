package utils

import (
	"encoding/json"
	"strings"

	"github.com/zigzter/chatterm/types"
)

func UsernoticeParser(input string) types.Message {
	parts := strings.SplitN(input, " :", 2)
	metadata := parts[0]
	message := ""
	messageSplit := strings.Split(parts[1], " :")
	if len(messageSplit) > 1 {
		message = messageSplit[1]
	}
	kvPairs := strings.Split(metadata, ";")
	tagMap := make(map[string]string)
	for _, kvPair := range kvPairs {
		kv := strings.Split(kvPair, "=")
		if len(kv) == 2 {
			tagMap[kv[0]] = kv[1]
		}
	}
	tagMap["message"] = message
	// This might be hacky. Dealing with the hyphenated keys of the IRC message tags
	jsonBytes, err := json.Marshal(tagMap)
	if err != nil {
	}
	var newMap types.MessageKVMap
	json.Unmarshal(jsonBytes, &newMap)
	switch tagMap["msg-id"] {
	case "announcement":
		return AnnouncementParser(newMap)
	case "sub":
		return GiftSubParser(newMap)
	case "resub":
		return SubParser(newMap)
	case "subgift":
	case "raid":
		return RaidParser(newMap)
	}
	return nil
}

func RaidParser(input types.MessageKVMap) types.RaidMessage {
	return types.RaidMessage{
		DisplayName: input.DisplayName,
		ViewerCount: input.ViewerCount,
	}
}

func AnnouncementParser(input types.MessageKVMap) types.AnnouncementMessage {
	return types.AnnouncementMessage{
		DisplayName: input.DisplayName,
		Message:     input.Message,
	}
}

func GiftSubParser(input types.MessageKVMap) types.SubGiftMessage {
	var subGiftMessage types.SubGiftMessage
	return subGiftMessage
}

func SubParser(input types.MessageKVMap) types.SubMessage {
	return types.SubMessage{
		Message:     input.Message,
		DisplayName: input.DisplayName,
		Months:      input.CumulativeMonths,
		Streak:      input.StreakMonths,
	}
}
