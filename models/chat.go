package models

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/twitch"
	"github.com/zigzter/chatterm/types"
	"github.com/zigzter/chatterm/utils"
)

type ChatModel struct {
	msgChan     chan types.ChatMessageWrap
	chatContent string
	input       string
	textinput   textinput.Model
	viewport    viewport.Model
	width       int
	height      int
	channel     string
	WsClient    *utils.WebSocketClient
}

func InitialChatModel(width int, height int) ChatModel {
	vp := viewport.New(width-2, height-2)
	log.Println("width: ", width, "height: ", height)
	vp.SetContent("")
	utils.InitConfig()
	username := viper.GetString("username")
	oauth := fmt.Sprintf("oauth:%s", viper.GetString("token"))
	channel := viper.GetString("channel")
	msgChan := make(chan types.ChatMessageWrap, 100)
	wsClient, err := utils.NewWebSocketClient()
	if err != nil {
		log.Fatal("Failed to initialize socket client")
	}
	go utils.EstablishWSConnection(wsClient, channel, username, oauth, msgChan)
	ti := textinput.New()
	ti.CharLimit = 256
	ti.Placeholder = "Send a message"
	ti.Focus()
	return ChatModel{
		input:     "",
		textinput: ti,
		viewport:  vp,
		msgChan:   msgChan,
		WsClient:  wsClient,
		channel:   channel,
	}
}

func processChatInput(input string) (isCommand bool, command string, args []string) {
	if strings.HasPrefix(input, "/") {
		parts := strings.SplitN(input, " ", 2)
		// The first part is the command
		command = strings.TrimPrefix(parts[0], "/")
		// The rest of the input (if any) is considered as arguments
		if len(parts) > 1 {
			args = strings.Split(parts[1], " ")
		}
		return true, command, args
	}
	return false, "", nil
}

func isValidCommand(command string) bool {
	switch types.TwitchCommand(command) {
	case types.Ban, types.Clear:
		return true
	}
	return false
}

func listenToWebSocket(msgChan <-chan types.ChatMessageWrap) tea.Cmd {
	return func() tea.Msg {
		return <-msgChan
	}
}

func (m ChatModel) Init() tea.Cmd {
	return listenToWebSocket(m.msgChan)
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			message := m.textinput.Value()
			isCommand, command, args := processChatInput(message)
			if isCommand {
				if isValidCommand(command) {
					// TODO: Give user feedback of result
					err := twitch.SendTwitchCommand(types.TwitchCommand(command), args)
					if err != nil {
						log.Println(err)
					}
				} else {
					m.chatContent += fmt.Sprintf("Invalid command: %s\n", command)
				}
			} else {
				m.chatContent += fmt.Sprintf("You: %s\n", message)
				m.WsClient.SendMessage([]byte("PRIVMSG #" + m.channel + " :" + message))
			}
			m.viewport.SetContent(m.chatContent)
			m.textinput.Reset()
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
		m.viewport.GotoBottom()
		return m, listenToWebSocket(m.msgChan)

	}
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m ChatModel) View() string {
	return fmt.Sprintf("%s\n%s", m.viewport.View(), m.textinput.View())
}
