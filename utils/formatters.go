package utils

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/zigzter/chatterm/types"
)

const (
	newMsgColor       = "#e64553"
	subColor          = "#04a5e5"
	announcementColor = "#40a02b"
	raidColor         = "#fe640b"
)

type boxWithLabel struct {
	BoxStyle   lipgloss.Style
	LabelStyle lipgloss.Style
}

func newBoxWithLabel(color string) boxWithLabel {
	return boxWithLabel{
		BoxStyle:   lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(color)).Padding(0),
		LabelStyle: lipgloss.NewStyle().Padding(0),
	}
}

func (b boxWithLabel) Render(label, content string) string {
	var (
		border          lipgloss.Border             = b.BoxStyle.GetBorderStyle()
		topBorderStyler func(strs ...string) string = lipgloss.NewStyle().Foreground(b.BoxStyle.GetBorderTopForeground()).Render
		topLeft         string                      = topBorderStyler(border.TopLeft)
		topRight        string                      = topBorderStyler(border.TopRight)
		renderedLabel   string                      = b.LabelStyle.Render(label)
	)
	width := lipgloss.Width(content)
	borderWidth := b.BoxStyle.GetHorizontalBorderSize()
	cellsShort := max(0, width+borderWidth-lipgloss.Width(topLeft+topRight+renderedLabel))
	gap := strings.Repeat(border.Top, cellsShort)
	top := topLeft + renderedLabel + topBorderStyler(gap) + topRight
	bottom := b.BoxStyle.Copy().
		BorderTop(false).
		Width(width).
		Render(content)
	return top + "\n" + bottom + "\n"
}

func usernameColorizer(color string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
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
	box := newBoxWithLabel(newMsgColor)
	msg := fmt.Sprintf(
		"[%s]%s%s: %s",
		message.Timestamp,
		icon,
		usernameColorizer(color).Render(message.DisplayName),
		message.Message,
	)
	if message.IsFirstMessage {
		return box.Render("First message", msg)
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
	box := newBoxWithLabel(subColor)
	msg := fmt.Sprintf(
		"%s subscribed for %s months%s",
		usernameColorizer(message.Color).Render(message.DisplayName),
		message.Months,
		fullMessage,
	)
	return box.Render("Sub", msg)
}

func FormatAnnouncementMessage(message types.AnnouncementMessage) string {
	box := newBoxWithLabel(announcementColor)
	msg := fmt.Sprintf(
		"[Announcement]%s: %s",
		usernameColorizer(message.Color).Render(message.DisplayName),
		message.Message,
	)
	return box.Render("Announcement", msg)
}

func FormatRaidMessage(message types.RaidMessage) string {
	box := newBoxWithLabel(raidColor)
	msg := fmt.Sprintf(
		"%s raided the channel with %s viewers!",
		usernameColorizer(message.Color).Render(message.DisplayName),
		message.ViewerCount,
	)
	return box.Render("Raid", msg)
}

func FormatGiftSubMessage(message types.SubGiftMessage) string {
	box := newBoxWithLabel(subColor)
	msg := fmt.Sprintf(
		"%s gifted a subscription to %s!",
		usernameColorizer(message.Color).Render(message.GiverName),
		message.ReceiverName,
	)
	return box.Render("Gift sub", msg)
}

func FormatMysteryGiftSubMessage(message types.MysterySubGiftMessage) string {
	box := newBoxWithLabel(subColor)
	msg := fmt.Sprintf(
		"%s is giving %s subs to the channel!",
		usernameColorizer(message.Color).Render(message.GiverName),
		message.GiftAmount,
	)
	return box.Render("Gifting subs", msg)
}
