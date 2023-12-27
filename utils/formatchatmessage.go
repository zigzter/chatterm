package utils

import (
	"fmt"

	"github.com/zigzter/chatterm/types"
)

func FormatChatMessage(message types.ChatMessage) string {
	icon := ""
	bgColor := ""
	if message.IsFirstMessage {
		bgColor = "\033[41m"
	}
	modIcon := "\033[32m󰓥"
	vipIcon := "\033[35m󰮊"
	if message.IsMod {
		icon = modIcon
	} else if message.IsVIP {
		icon = vipIcon
	}

	color := ParseHexColor(message.Color)
	resetCode := "\033[0m"
	defaultTextColor := "\033[39m"
	return fmt.Sprintf(
		"[%s]%s%s\033[38;2;%d;%d;%dm%s%s%s: %s%s\n",
		message.Timestamp,
		bgColor,
		icon,
		color.R, color.G, color.B,
		message.DisplayName,
		defaultTextColor,
		resetCode,
		message.Message,
		resetCode,
	)
}
