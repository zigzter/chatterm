package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zigzter/chatterm/twitch"
	"github.com/zigzter/chatterm/types"
	"github.com/zigzter/chatterm/utils"
)

const (
	minUsernameLength = 4
	maxUsernameLength = 24
)

func isValidUsernameLength(input textinput.Model) bool {
	length := len(input.Value())
	return length >= minUsernameLength && length <= maxUsernameLength
}

type AuthModel struct {
	input             textinput.Model
	spinner           spinner.Model
	serverStarting    bool
	serverStarted     bool
	authPromptOpening bool
	authPrompOpened   bool
	tokenReceiving    bool
	tokenReceived     bool
	externalMsgs      chan tea.Msg
}

func InitialAuthModel() AuthModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	ti := textinput.New()
	ti.Placeholder = "Twitch username"
	ti.CharLimit = maxUsernameLength
	ti.Focus()
	return AuthModel{
		input:             ti,
		spinner:           s,
		externalMsgs:      make(chan tea.Msg, 10),
		serverStarting:    false,
		serverStarted:     false,
		authPromptOpening: false,
		authPrompOpened:   false,
		tokenReceiving:    false,
		tokenReceived:     false,
	}
}

func listenForExternalMsgs(externalMsgs chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-externalMsgs
	}
}

func (m AuthModel) Init() tea.Cmd {
	return listenForExternalMsgs(m.externalMsgs)
}

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			if isValidUsernameLength(m.input) {
				username := m.input.Value()
				utils.SaveConfig(map[string]interface{}{
					utils.UsernameKey: username,
				})
				m.serverStarting = true
				ready := make(chan struct{}, 1)
				cmds = append(
					cmds,
					twitch.StartLocalServer(ready, m.externalMsgs),
					listenForExternalMsgs(m.externalMsgs),
				)
			}
		case tea.KeyCtrlI:
			return ChangeView(m, ChannelInputState)
		}
	case types.ServerStartedMsg:
		m.serverStarting = false
		m.serverStarted = true
		cmds = append(
			cmds,
			twitch.PromptTwitchAuth(),
			listenForExternalMsgs(m.externalMsgs),
		)
	case types.AuthOpenMsg:
		m.authPromptOpening = true
	case types.AuthOpenedMsg:
		m.authPromptOpening = false
		m.authPrompOpened = true
	case types.TokenReceiveMsg:
		m.tokenReceiving = true
		m.authPrompOpened = false
	case types.TokenReceivedMsg:
		m.tokenReceiving = false
		m.tokenReceived = true
	}
	if cmd := listenForExternalMsgs(m.externalMsgs); cmd != nil {
		cmds = append(cmds, cmd)
	}
	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	if inputCmd != nil {
		cmds = append(cmds, inputCmd)
	}
	return m, tea.Batch(cmds...)
}

func (m AuthModel) View() string {
	var b strings.Builder
	check := ""
	fmt.Fprintf(
		&b,
		"Enter username:\n%s\n",
		m.input.View(),
	)
	if isValidUsernameLength(m.input) {
		fmt.Fprintln(&b, "Press [enter] to start authentication process")
	}
	if m.serverStarting {
		fmt.Fprintln(&b, "Starting server...", m.spinner.View())
	} else if m.serverStarted {
		fmt.Fprintln(&b, "Starting server...", check)
	}
	if m.authPromptOpening {
		fmt.Fprintln(&b, "Opening Twitch authentication...", m.spinner.View())
	} else if m.serverStarted {
		fmt.Fprintln(&b, "Opening Twitch authentication...", check)
	}
	if m.tokenReceiving {
		fmt.Fprintln(&b, "Processing auth token...", m.spinner.View())
	} else if m.tokenReceived {
		fmt.Fprintln(&b, "Processing auth token...", check)
		fmt.Fprintln(&b, "Authentication complete! Press [Ctrl+i] to return to channel selection")
	}
	return b.String()
}
