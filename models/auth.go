package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zigzter/chatterm/twitch"
	"github.com/zigzter/chatterm/types"
)

type AuthModel struct {
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
	return AuthModel{
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
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "f":
			m.serverStarting = true
			ready := make(chan struct{}, 1)
			cmd = tea.Batch(
				twitch.StartLocalServer(ready, m.externalMsgs),
				listenForExternalMsgs(m.externalMsgs),
			)
		case "c":
			return ChangeView(m, ChannelInputState)
		}
	case types.ServerStartedMsg:
		m.serverStarting = false
		m.serverStarted = true
		cmd = tea.Batch(
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
	if cmd == nil {
		cmd = listenForExternalMsgs(m.externalMsgs)
	}
	return m, cmd
}

func (m AuthModel) View() string {
	var b strings.Builder
	check := ""
	fmt.Fprint(&b, "Press [f] to start auth, [c] to cancel\n")
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
		fmt.Fprintln(&b, "Authentication complete! Press [c] to return to channel selection")
	}
	return b.String()
}
