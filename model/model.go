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
	messages    []types.ChatMessage
	msgChan     chan types.ChatMessageWrap
	chatContent string
	input       string
	textinput   textinput.Model
	viewport    viewport.Model
	width       int
	height      int
	WsClient    *utils.WebSocketClient
}

func InitialModel() model {
	vp := viewport.New(84, 24)
	vp.SetContent("")
	utils.InitConfig()
	username := viper.GetString("username")
	oauth := viper.GetString("oauth")
	msgChan := make(chan types.ChatMessageWrap, 100)
	wsClient, err := utils.NewWebSocketClient()
	if err != nil {
		log.Fatal("Failed to initialize socket client")
	}
	go utils.EstablishWSConnection(wsClient, "a_seagull", username, oauth, msgChan)
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Focus()
	return model{
		messages:  []types.ChatMessage{},
		input:     "",
		textinput: ti,
		viewport:  vp,
		msgChan:   msgChan,
		WsClient:  wsClient,
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
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.chatContent += fmt.Sprintf("You: %s\n", m.textinput.Value())
			m.textinput.SetValue("")
			return m, listenToWebSocket(m.msgChan)
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.viewport.Width = m.width
		m.viewport.Height = m.height
		return m, listenToWebSocket(m.msgChan)
	case types.ChatMessageWrap:
		switch msg := msg.ChatMsg.(type) {
		case types.ChatMessage:
			m.chatContent += utils.FormatChatMessage(msg)
			m.viewport.SetContent(m.chatContent)
		case types.SubMessage:
			m.chatContent += utils.FormatSubMessage(msg)
			m.viewport.SetContent(m.chatContent)
		}
		m.viewport.YOffset = m.viewport.TotalLineCount() - m.viewport.Height
		if m.viewport.YOffset < 0 {
			m.viewport.YOffset = 0
		}
		return m, listenToWebSocket(m.msgChan)

	}
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf("%s\n%s", m.viewport.View(), m.textinput.View())
}
