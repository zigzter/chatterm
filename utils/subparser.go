package utils

import (
    "strings"
    "fmt"
)

type SubMessage struct {
    DisplayName string
    Message string
    Months string
    Streak string
}

func SubParser(input string) SubMessage {
	parts := strings.SplitN(input, " :", 2)
	metadata := parts[0]
	message := strings.Split(parts[1], " :")[1]
	keyValPairs := strings.Split(metadata, ";")
    subMessage := SubMessage{ Message: message }
    for _, kvPair := range keyValPairs {
        kv := strings.Split(kvPair, "=")
        if len(kv) == 2 {
            key := kv[0]
            value := kv[1]
            fmt.Println(key, value)
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
