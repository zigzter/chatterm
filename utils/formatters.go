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

func highlighter(msg string, color string) string {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(color)).
		BorderTop(true).
		BorderRight(true).
		BorderBottom(true).
		Padding(0).
		Margin(0).
		Render(msg) + "\n"
}

func iconColorizer(color string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func FormatChatMessage(message types.ChatMessage) string {
	icon := ""
	modIcon := iconColorizer("#40a02b").Render("[󰓥]")
	vipIcon := iconColorizer("#ea76cb").Render("[󰮊]")
	if message.IsMod {
		icon = modIcon
	} else if message.IsVIP {
		icon = vipIcon
	}
	color := message.Color
	if color == "" {
		color = "#7287fd"
	}
	msg := fmt.Sprintf(
		"[%s]%s%s: %s",
		message.Timestamp,
		icon,
		usernameColorizer(color).Render(message.DisplayName),
		message.Message,
	)
	if message.IsFirstMessage {
		return highlighter(msg, "#e64553")
	} else {
		return msg + "\n"
	}
}

func FormatSubMessage(message types.SubMessage) string {
	var fullMessage string
	if message.Message != "" {
		fullMessage = ": " + message.Message
	} else {
		fullMessage = "!"
	}
	msg := fmt.Sprintf(
		"[%s]%s subscribed for %s months%s",
		message.Timestamp,
		usernameColorizer(message.Color).Render(message.DisplayName),
		message.Months,
		fullMessage,
	)
	return highlighter(msg, "12")
}

func FormatAnnouncementMessage(message types.AnnouncementMessage) string {
	msg := fmt.Sprintf(
		"[%s][Announcement]%s: %s",
		message.Timestamp,
		usernameColorizer(message.Color).Render(message.DisplayName),
		message.Message,
	)
	return highlighter(msg, "228")
}

func FormatRaidMessage(message types.RaidMessage) string {
	msg := fmt.Sprintf(
		"[%s]%s raided the channel with %s viewers!",
		message.Timestamp,
		usernameColorizer(message.Color).Render(message.DisplayName),
		message.ViewerCount,
	)
	return highlighter(msg, "12")
}

func FormatGiftSubMessage(message types.SubGiftMessage) string {
	msg := fmt.Sprintf(
		"[%s]%s gifted a subscription to %s",
		message.Timestamp,
		usernameColorizer(message.Color).Render(message.GiverName),
		message.ReceiverName,
	)
	return highlighter(msg, "12")
}

func FormatMysteryGiftSubMessage(message types.MysterySubGiftMessage) string {
	msg := fmt.Sprintf(
		"[%s]%s is giving %s subs to the channel!",
		message.Timestamp,
		usernameColorizer(message.Color).Render(message.GiverName),
		message.GiftAmount,
	)
	return highlighter(msg, "9")
}
