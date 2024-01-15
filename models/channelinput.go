package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/twitch"
	"github.com/zigzter/chatterm/types"
	"github.com/zigzter/chatterm/utils"
)

var (
	titleStyle       = lipgloss.NewStyle()
	channelNameStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	viewerCountStyle = lipgloss.NewStyle()
	gameNameStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
)

type ChannelInputModel struct {
	textinput    textinput.Model
	authRequired bool
	liveStreams  []types.LiveChannelsData
	error        string
}

func InitialChannelInputModel() ChannelInputModel {
	utils.InitConfig()
	authRequired := utils.IsAuthRequired()
	var error string
	userID := viper.GetString("userid")
	if userID == "" {
		username := viper.GetString("username")
		userData, err := twitch.SendUserRequest(username)
		if err != nil {
			error = err.Error()
		} else {
			userID = userData.Data[0].ID
			utils.SaveConfig(map[string]interface{}{
				"userid": userID,
			})
		}
	}
	var liveStreams []types.LiveChannelsData
	liveStreamsResp, err := twitch.SendLiveChannelsRequest(userID)
	if err != nil {
		error = err.Error()
	} else {
		liveStreams = liveStreamsResp.Data
		error = ""
	}
	ti := textinput.New()
	ti.Placeholder = "a_seagull"
	configChannel := viper.GetString("channel")
	if configChannel != "" {
		ti.SetValue(configChannel)
	}
	ti.Focus()
	return ChannelInputModel{
		textinput:    ti,
		authRequired: authRequired,
		liveStreams:  liveStreams,
		error:        error,
	}
}

func (m ChannelInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ChannelInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
			// TODO: re-enable config view
			// case tea.KeyCtrlO:
			// 	return ChangeView(m, ConfigState)
		case tea.KeyCtrlR:
			userID := viper.GetString("userid")
			if userID == "" {
				username := viper.GetString("username")
				userData, err := twitch.SendUserRequest(username)
				if err != nil {
					m.error = err.Error()
				}
				userID = userData.Data[0].ID
				utils.SaveConfig(map[string]interface{}{
					"userid": userID,
				})
			}
			liveStreamsResp, err := twitch.SendLiveChannelsRequest(userID)
			if err != nil {
				m.error = err.Error()
			} else {
				m.liveStreams = liveStreamsResp.Data
				m.error = ""
			}
		case tea.KeyCtrlA:
			return ChangeView(m, AuthState)
		case tea.KeyEnter:
			if m.textinput.Value() == "exit" {
				return m, tea.Quit
			}
			utils.SaveConfig(map[string]interface{}{
				"channel": m.textinput.Value(),
			})
			return ChangeView(m, ChatState)
		}
	}
	m.authRequired = utils.IsAuthRequired()
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m ChannelInputModel) View() string {
	var b strings.Builder
	if m.authRequired {
		b.WriteString("Authentication required. Press [Ctrl+a] to start.")
	} else {
		b.WriteString("Live Channels:\n")
		if m.error != "" {
			b.WriteString("Live channel retrieval error:" + m.error)
		}
		for i, channel := range m.liveStreams {
			if i > 10 {
				break
			}
			b.WriteString(fmt.Sprintf(
				"%s playing %s: %s (%s viewers)\n",
				channelNameStyle.Render(channel.UserName),
				gameNameStyle.Render(channel.GameName),
				titleStyle.Render(channel.Title),
				viewerCountStyle.Render(fmt.Sprintf("%d", channel.ViewerCount)),
			))
		}
		b.WriteString("\nEnter channel name:\n")
		b.WriteString(m.textinput.View() + "\n")
	}
	b.WriteString(helpStyle.Render("[Ctrl+c]: quit - [Ctrl+r]: reload streams"))
	return b.String()
}
