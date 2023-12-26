package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/zigzter/chatterm/types"
)

type ChatMessageWrap struct {
	ChatMsg types.ChatMessage
}

type model struct {
	messages     []types.ChatMessage
	input        string
	inputFocused bool
	textinput    textinput.Model
	viewport     viewport.Model
}

func InitialModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 256
	vp := viewport.New(30, 5)
	vp.SetContent("Hello there...")
	return model{
		messages:     []types.ChatMessage{},
		input:        "",
		inputFocused: true,
		textinput:    ti,
		viewport:     vp,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

type (
	inputMessage   string
	newChatMessage string
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			m.input += string(msg.Runes)
		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case tea.KeyEnter:
			// TODO: send to Twitch here
			m.textinput.Update("")
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case ChatMessageWrap:
		m.messages = append(m.messages, msg.ChatMsg)
	}
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textinput.View(),
	) + "\n\n"
}
