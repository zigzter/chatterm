package utils

import (
    "fmt"
)

func PrintSubMessage(message SubMessage) {
    fmt.Printf("%s subscribed for %s months: %s", message.DisplayName, message.Months, message.Message)
}
