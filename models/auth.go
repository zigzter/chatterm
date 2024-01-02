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
	fetching bool
	error    string
	spinner  spinner.Model
}

func InitialAuthModel() AuthModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return AuthModel{
		spinner:  s,
		fetching: false,
		error:    "",
	}
}

func (m AuthModel) Init() tea.Cmd {
	return nil
}

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "f":
			m.fetching = true
			return m, twitch.StartAuthenticationProcess()
		case "c":
			return ChangeView(m, ChannelInputState)
		}
	case types.AuthResultMsg:
		m.fetching = false
		m.error = msg.Error
		return m, nil
	}
	return m, nil
}

func (m AuthModel) View() string {
	var b strings.Builder
	if m.fetching {
		fmt.Fprintf(&b, "Fetching token... %s", m.spinner.View())
	} else {
		fmt.Fprintln(&b, "Press 'f' to fetch token, 'c' to cancel.")
	}
	if m.error != "" && !m.fetching {
		fmt.Fprintln(&b, "Error:", m.error)
	}
	if m.error == "" && !m.fetching {
		fmt.Fprintln(&b, "Authentication successful!")
	}
	return b.String()
}
