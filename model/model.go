package model

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/types"
	"github.com/zigzter/chatterm/utils"
)

type model struct {
	messages    []types.ChatMessage
	msgChan     chan types.ChatMessageWrap
	chatContent string
	input       string
	textinput   textinput.Model
	viewport    viewport.Model
	width       int
	height      int
}

func InitialModel() model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 256
	vp := viewport.New(84, 24)
	vp.SetContent("")
	utils.InitConfig()
	username := viper.GetString("username")
	oauth := viper.GetString("oauth")
	msgChan := make(chan types.ChatMessageWrap, 100)
	go utils.EstablishWSConnection("summit1g", username, oauth, msgChan)
	return model{
		messages:  []types.ChatMessage{},
		input:     "",
		textinput: ti,
		viewport:  vp,
		msgChan:   msgChan,
	}
}

func listenToWebSocket(msgChan <-chan types.ChatMessageWrap) tea.Cmd {
	return func() tea.Msg {
		return <-msgChan
	}
}

func (m model) Init() tea.Cmd {
	return listenToWebSocket(m.msgChan)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.viewport.Width = m.width
		m.viewport.Height = m.height
		return m, listenToWebSocket(m.msgChan)
	case types.ChatMessageWrap:
		m.chatContent += utils.FormatChatMessage(msg.ChatMsg)
		m.viewport.SetContent(m.chatContent)
		return m, listenToWebSocket(m.msgChan)
	}
	return m, listenToWebSocket(m.msgChan)
}

func (m model) View() string {
	return m.viewport.View()
}
