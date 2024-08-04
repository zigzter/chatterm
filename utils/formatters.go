package utils

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/types"
)

const (
	subColor          = "#04a5e5"
	announcementColor = "#40a02b"
	raidColor         = "#fe640b"
)

type BoxWithLabel struct {
	BoxStyle   lipgloss.Style
	LabelStyle lipgloss.Style
	width      int
}

func NewBoxWithLabel(color string) BoxWithLabel {
	return BoxWithLabel{
		BoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(color)).
			Padding(0),
		LabelStyle: lipgloss.NewStyle().Padding(0),
	}
}

// SetWidth is used for wrapping contained text
func (b *BoxWithLabel) SetWidth(width int) *BoxWithLabel {
	b.width = width
	return b
}

func (b *BoxWithLabel) Render(label, content string) string {
	var (
		topBorderStyler func(strs ...string) string = lipgloss.NewStyle().
				Foreground(b.BoxStyle.GetBorderTopForeground()).
				Render
		border        lipgloss.Border = b.BoxStyle.GetBorderStyle()
		topLeft       string          = topBorderStyler(border.TopLeft)
		topRight      string          = topBorderStyler(border.TopRight)
		renderedLabel string          = b.LabelStyle.Render(label)
	)
	width := lipgloss.Width(content)
	if b.width != 0 {
		width = b.width
	}
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

var (
	showBadges            bool = true
	showTimestamps        bool = true
	highlightSubs         bool
	highlightRaids        bool
	firstTimeChatterColor string
	watchedUsers          map[string]any
)

// SetFormatterConfigValues sets the formatter customization options from the config.
// This is required because Viper won't have loaded the config yet when it reads this file.
func SetFormatterConfigValues() {
	showBadges = viper.GetBool(ShowBadgesKey)
	showTimestamps = viper.GetBool(ShowTimestampsKey)
	highlightSubs = viper.GetBool(HighlightRaidsKey)
	highlightRaids = viper.GetBool(HighlightRaidsKey)
	watchedUsers = viper.GetStringMap(WatchedUsersKey)
	firstTimeChatterColor = viper.GetString(FirstTimeChatterColorKey)
}

// GenerateIcon returns a colored user-type icon, if applicable to the user.
// For example, a green sword icon for a moderator.
func GenerateIcon(userType string) string {
	if !showBadges {
		return ""
	}
	switch userType {
	case "broadcaster":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#d20f39")).Render("[]")
	case "moderator":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#40a02b")).Render("[󰓥]")
	case "vip":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ea76cb")).Render("[󰮊]")
	case "staff":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#8839ef")).Render("[󰖷]")
	}
	return ""
}

func FormatChatMessage(message types.ChatMessage, width int) string {
	icon := GenerateIcon(message.ChannelUserType)
	color := message.Color
	if color == "" {
		color = "#7287fd"
	}
	timestamp := ""
	if showTimestamps {
		timestamp = "[" + ParseTimestamp(message.Timestamp).Time + "]"
	}
	box := NewBoxWithLabel(firstTimeChatterColor)
	msg := fmt.Sprintf(
		"%s%s%s: %s",
		timestamp,
		icon,
		usernameColorizer(color).Render(message.DisplayName),
		message.Message,
	)
	msg = wordwrap.String(msg, width)
	if message.IsFirstMessage {
		return box.Render("First message", msg)
	} else if watchedUsers[strings.ToLower(message.DisplayName)] == true {
		return box.Render("Watched user", msg)
	} else {
		return msg + "\n"
	}
}

func FormatSubMessage(message types.SubMessage, width int) string {
	var fullMessage string
	if message.Message != "" {
		fullMessage = ": " + message.Message
	} else {
		fullMessage = "!"
	}
	box := NewBoxWithLabel(subColor)
	msg := fmt.Sprintf(
		"%s subscribed for %s months%s",
		usernameColorizer(message.Color).Render(message.DisplayName),
		message.Months,
		fullMessage,
	)
	msg = wordwrap.String(msg, width)
	if highlightSubs {
		return box.Render("Sub", msg)
	}
	return msg + "\n"
}

func FormatAnnouncementMessage(message types.AnnouncementMessage, width int) string {
	box := NewBoxWithLabel(announcementColor)
	msg := fmt.Sprintf(
		"%s: %s",
		usernameColorizer(message.Color).Render(message.DisplayName),
		message.Message,
	)
	msg = wordwrap.String(msg, width)
	return box.Render("Announcement", msg)
}

func FormatRaidMessage(message types.RaidMessage, width int) string {
	box := NewBoxWithLabel(raidColor)
	msg := fmt.Sprintf(
		"%s raided the channel with %s viewers!",
		usernameColorizer(message.Color).Render(message.DisplayName),
		message.ViewerCount,
	)
	msg = wordwrap.String(msg, width)
	if highlightRaids {
		return box.Render("Raid", msg)
	}
	return msg + "\n"
}

func FormatGiftSubMessage(message types.SubGiftMessage, width int) string {
	box := NewBoxWithLabel(subColor)
	msg := fmt.Sprintf(
		"%s gifted a subscription to %s!",
		usernameColorizer(message.Color).Render(message.GiverName),
		message.ReceiverName,
	)
	msg = wordwrap.String(msg, width)
	if highlightSubs {
		return box.Render("Gift sub", msg)
	}
	return msg + "\n"
}

func FormatMysteryGiftSubMessage(message types.MysterySubGiftMessage, width int) string {
	box := NewBoxWithLabel(subColor)
	msg := fmt.Sprintf(
		"%s is giving %s subs to the channel!",
		usernameColorizer(message.Color).Render(message.GiverName),
		message.GiftAmount,
	)
	msg = wordwrap.String(msg, width)
	if highlightSubs {
		return box.Render("Gifting subs", msg)
	}
	return msg + "\n"
}
