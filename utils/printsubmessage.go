package utils

import (
	"fmt"

	"github.com/zigzter/chatterm/types"
)

func PrintSubMessage(message types.SubMessage) {
	fmt.Printf(
		"%s subscribed for %s months: %s",
		message.DisplayName,
		message.Months,
		message.Message,
	)
}
