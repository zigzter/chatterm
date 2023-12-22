package utils

import (
	"strings"
)

type ChatMessage struct {
	Timestamp      string
    // comma seperated string, each badge followed by a slash and a number
	Badges         string
    // User color in hex
	Color          string
	DisplayName    string
	IsFirstMessage bool
	IsMod          bool
	IsVIP          bool
	Message        string
}

func MessageParser(input string) ChatMessage {
	parts := strings.SplitN(input, " :", 2)
	metadata := parts[0]
	message := strings.Split(parts[1], " :")[1]
	keyValPairs := strings.Split(metadata, ";")
	chatMessage := ChatMessage{Message: message}
	for _, kvPair := range keyValPairs {
		kv := strings.Split(kvPair, "=")
		if len(kv) == 2 {
			key := kv[0]
			value := kv[1]

			switch key {
			case "badges":
				chatMessage.Badges = value
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
			}
		}
	}

	// Parse the timestamp and username
	// parts = strings.SplitN(parts[0], " ", 4)
	// chatMessage.Timestamp = parts[0] + " " + parts[1]
	// chatMessage.Username = parts[3]
	return chatMessage
}
