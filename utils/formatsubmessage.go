package utils

import (
	"fmt"

	"github.com/zigzter/chatterm/types"
)

func FormatSubMessage(message types.SubMessage) string {
	return fmt.Sprintf(
		"%s subscribed for %s months: %s\n",
		message.DisplayName,
		message.Months,
		message.Message,
	)
}
