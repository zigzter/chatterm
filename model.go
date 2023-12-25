package main

import (
	"github.com/charmbracelet/bubbletea"
)

type model struct {
	messages     []string
	input        string
	inputFocused bool
}

func initialModel() model {
	return model{
		messages:     []string{},
		input:        "",
		inputFocused: true,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

type (
	inputMessage   string
	newChatMessage string
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyRunes {
			m.input += string(msg.Runes)
		} else if msg.Type == tea.KeyBackspace {
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		} else if msg.Type == tea.KeyEnter {
			// TODO: send to Twitch here
			m.input = ""
		}
	case newChatMessage:
		m.messages = append(m.messages, string(msg))
	}
	return m, nil
}

func (m model) View() string {
	return "Hello"
}
