package utils

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/zigzter/chatterm/types"
)

func usernameColorizer(color string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func firstMessageColorizer() lipgloss.Style {
	return lipgloss.NewStyle().Background(lipgloss.Color("#fe640b"))
}

func FormatChatMessage(message types.ChatMessage) string {
	icon := ""
	msg := message.Message
	if message.IsFirstMessage {
		msg = firstMessageColorizer().Render(msg)
	}
	modIcon := "\033[32m󰓥"
	vipIcon := "\033[35m󰮊"
	if message.IsMod {
		icon = modIcon
	} else if message.IsVIP {
		icon = vipIcon
	}
	return fmt.Sprintf(
		"[%s]%s%s: %s\n",
		message.Timestamp,
		icon,
		usernameColorizer(message.Color).Render(message.DisplayName),
		msg,
	)
}

func FormatSubMessage(message types.SubMessage) string {
	var fullMessage string
	if message.Message != "" {
		fullMessage = ": " + message.Message
	} else {
		fullMessage = "!"
	}
	return fmt.Sprintf(
		"[%s]%s subscribed for %s months%s\n",
		message.Timestamp,
		usernameColorizer(message.Color).Render(message.DisplayName),
		message.Months,
		fullMessage,
	)
}

func FormatAnnouncementMessage(message types.AnnouncementMessage) string {
	return fmt.Sprintf(
		"[%s][Announcement]%s: %s\n",
		message.Timestamp,
		usernameColorizer(message.Color).Render(message.DisplayName),
		message.Message,
	)
}

func FormatRaidMessage(message types.RaidMessage) string {
	return fmt.Sprintf(
		"[%s]%s raided the channel with %s viewers!\n",
		message.Timestamp,
		usernameColorizer(message.Color).Render(message.DisplayName),
		message.ViewerCount,
	)
}

func FormatGiftSubMessage(message types.SubGiftMessage) string {
	return fmt.Sprintf(
		"[%s]%s gifted a subscription to %s\n",
		message.Timestamp,
		usernameColorizer(message.Color).Render(message.GiverName),
		message.ReceiverName,
	)
}

func FormatMysteryGiftSubMessage(message types.MysterySubGiftMessage) string {
	return fmt.Sprintf(
		"[%s]%s is giving %s subs to the channel!\n",
		message.Timestamp,
		usernameColorizer(message.Color).Render(message.GiverName),
		message.GiftAmount,
	)
}
