package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
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

type KeyMap struct {
	Auth    key.Binding
	Options key.Binding
	Tab     key.Binding
	Input   key.Binding
	Refresh key.Binding
	Esc     key.Binding
	Close   key.Binding
	Enter   key.Binding
	Prev    key.Binding
	Next    key.Binding
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Prev, k.Next}, {k.Input, k.Refresh}, {k.Auth, k.Close}, {k.Options},
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Enter, k.Esc, k.Close}
}

var keys = KeyMap{
	Auth:    key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "open auth")),
	Options: key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "open options")),
	Tab:     key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "autocomplete")),
	Input:   key.NewBinding(key.WithKeys("i"), key.WithHelp("i", "open channel input")),
	Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh channels")),
	Close:   key.NewBinding(key.WithKeys("ctrl+c"), key.WithHelp("ctrl+c", "close app")),
	Esc:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "close channel input")),
	Enter:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "join channel")),
	Prev:    key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "previous page")),
	Next:    key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "next page")),
}

type ChannelInputModel struct {
	textinput         textinput.Model
	authRequired      bool
	liveStreams       []types.LiveChannelsData
	error             string
	ac                *utils.Trie
	width             int
	height            int
	paginator         paginator.Model
	streamsViewHeight int
	inputVisible      bool
	help              help.Model
}

func InitialChannelInputModel() ChannelInputModel {
	model := ChannelInputModel{
		width:        0,
		inputVisible: false,
		help:         help.New(),
	}
	model.help.ShowAll = true
	model.ac = &utils.Trie{Root: utils.NewTrieNode()}
	utils.InitConfig()
	authRequired := utils.IsAuthRequired()
	ti := textinput.New()
	ti.Placeholder = "a_seagull"
	ti.Blur()
	configChannel := viper.GetString(utils.ChannelKey)
	if configChannel != "" {
		ti.SetValue(configChannel)
	}
	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 10
	p.ActiveDot = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).
		Render("•")
	p.InactiveDot = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).
		Render("•")
	model.paginator = p
	model.textinput = ti
	model.authRequired = authRequired
	fetchLiveStreams(&model)
	return model
}

func fetchLiveStreams(m *ChannelInputModel) {
	userID := viper.GetString(utils.UserIDKey)
	if userID == "" {
		username := viper.GetString(utils.UsernameKey)
		userData, err := twitch.SendUserRequest(username)
		if err != nil {
			m.error = err.Error()
		} else {
			userID = userData.Data[0].ID
			utils.SaveConfig(map[string]interface{}{
				utils.UserIDKey: userID,
			})
		}
	}
	liveStreamsResp, err := twitch.SendLiveChannelsRequest(userID)
	if err != nil {
		m.error = err.Error()
	} else {
		m.liveStreams = liveStreamsResp.Data
		m.paginator.SetTotalPages(len(m.liveStreams))
		var streamNames []string
		for _, stream := range liveStreamsResp.Data {
			streamNames = append(streamNames, stream.UserName)
		}
		m.ac.Populate(streamNames)
		m.error = ""
	}
}

func (m ChannelInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ChannelInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Next):
			// Manually firing these to be able to disable them when the input is enabled
			if !m.inputVisible {
				m.paginator.NextPage()
			}
		case key.Matches(msg, keys.Prev):
			if !m.inputVisible {
				m.paginator.PrevPage()
			}
		case key.Matches(msg, keys.Input):
			if !m.inputVisible {
				m.inputVisible = true
				m.textinput.Focus()
				m.help.ShowAll = false
				return m, nil
			}
		case key.Matches(msg, keys.Esc):
			m.help.ShowAll = true
			m.inputVisible = false
			m.textinput.Blur()
			return m, nil
		case key.Matches(msg, keys.Tab):
			input := m.textinput.Value()
			suggestion := m.ac.UpdateSuggestion(input)
			m.textinput.SetValue(suggestion)
			m.textinput.CursorEnd()
		case key.Matches(msg, keys.Options):
			if !m.inputVisible {
				return ChangeView(m, SettingsState)
			}
		case key.Matches(msg, keys.Auth):
			if !m.inputVisible {
				return ChangeView(m, AuthState)
			}
		case key.Matches(msg, keys.Refresh):
			if !m.inputVisible {
				fetchLiveStreams(&m)
				return m, nil
			}
		case key.Matches(msg, keys.Close):
			return m, tea.Quit
		case key.Matches(msg, keys.Enter):
			utils.SaveConfig(map[string]interface{}{
				utils.ChannelKey: m.textinput.Value(),
			})
			return ChangeView(m, ChatState)
		}
	}
	m.authRequired = utils.IsAuthRequired()
	var tiCmd, pCmd tea.Cmd
	m.textinput, tiCmd = m.textinput.Update(msg)
	if !m.inputVisible {
		m.paginator, pCmd = m.paginator.Update(msg)
	}
	return m, tea.Batch(tiCmd, pCmd)
}

func (m ChannelInputModel) View() string {
	var b strings.Builder
	if m.authRequired {
		b.WriteString("Authentication required. Press [a] to start.\n")
	} else {
		b.WriteString("Live Channels:\n")
		if m.error != "" {
			b.WriteString("Live channel retrieval error:" + m.error)
		}
		start, end := m.paginator.GetSliceBounds(len(m.liveStreams))
		totalHeight := 0
		for _, channel := range m.liveStreams[start:end] {
			liveChannel := fmt.Sprintf(
				"%s playing %s: %s (%s viewers)\n",
				channelNameStyle.Render(channel.UserName),
				gameNameStyle.Render(channel.GameName),
				titleStyle.Render(channel.Title),
				viewerCountStyle.Render(fmt.Sprintf("%d", channel.ViewerCount)),
			)
			wrapped := wordwrap.String(liveChannel, m.width-4)
			// Single lines return 2 from Height, so we're subtracting 1
			totalHeight += lipgloss.Height(wrapped) - 1
			b.WriteString(wrapped)
		}
		m.streamsViewHeight = totalHeight
		b.WriteString("  " + m.paginator.View() + "\n")
		if m.inputVisible {
			b.WriteString("\nEnter channel name:\n")
			b.WriteString(m.textinput.View() + "\n")
		}
	}
	b.WriteString(m.help.View(keys))
	return b.String()
}
