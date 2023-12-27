package model

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/types"
	"github.com/zigzter/chatterm/utils"
)

type model struct {
	messages     []types.ChatMessage
	msgChan      chan types.ChatMessageWrap
	chatContent  string
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
	vp.SetContent("")
	return model{
		messages:     []types.ChatMessage{},
		input:        "",
		inputFocused: true,
		textinput:    ti,
		viewport:     vp,
	}
}

func (m model) Init() tea.Cmd {
	utils.InitConfig()
	username := viper.GetString("username")
	oauth := viper.GetString("oauth")
	m.msgChan = make(chan types.ChatMessageWrap, 50)
	go utils.EstablishWSConnection("Flats", username, oauth, m.msgChan)
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("Received message type: %T\n", msg)
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
	case types.ChatMessageWrap:
		m.chatContent += fmt.Sprintf("%s: %s\n", msg.ChatMsg.DisplayName, msg.ChatMsg.Message)
		m.viewport.SetContent(m.chatContent)
	default:
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textinput.View(),
	) + "\n\n"
}
