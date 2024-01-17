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
		// TODO: handle json error
	}
	var newMap types.MessageKVMap
	json.Unmarshal(jsonBytes, &newMap)
	newMap.Timestamp = ParseTimestamp(tagMap["tmi-sent-ts"])
	switch tagMap["msg-id"] {
	case "announcement":
		return AnnouncementParser(newMap)
	case "sub":
		return SubParser(newMap)
	case "resub":
		return SubParser(newMap)
	case "subgift":
		return GiftSubParser(newMap)
	case "submysterygift": // Gifting n subs to the channel
		return MysteryGiftSubParser(newMap)
	case "raid":
		return RaidParser(newMap)
	}
	return nil
}

func RaidParser(input types.MessageKVMap) types.RaidMessage {
	return types.RaidMessage{
		DisplayName: input.DisplayName,
		ViewerCount: input.ViewerCount,
		Timestamp:   input.Timestamp,
	}
}

func AnnouncementParser(input types.MessageKVMap) types.AnnouncementMessage {
	return types.AnnouncementMessage{
		DisplayName: input.DisplayName,
		Message:     input.Message,
		Timestamp:   input.Timestamp,
	}
}

func GiftSubParser(input types.MessageKVMap) types.SubGiftMessage {
	return types.SubGiftMessage{
		GiverName:    input.DisplayName,
		ReceiverName: input.ReceiverName,
		Timestamp:    input.Timestamp,
	}
}

func MysteryGiftSubParser(input types.MessageKVMap) types.MysterySubGiftMessage {
	return types.MysterySubGiftMessage{
		GiverName:  input.DisplayName,
		GiftAmount: input.GiftAmount,
		Timestamp:  input.Timestamp,
	}
}

func SubParser(input types.MessageKVMap) types.SubMessage {
	return types.SubMessage{
		Message:     input.Message,
		DisplayName: input.DisplayName,
		Months:      input.CumulativeMonths,
		Streak:      input.StreakMonths,
		Timestamp:   input.Timestamp,
	}
}
