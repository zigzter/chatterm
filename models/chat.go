package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/twitch"
	"github.com/zigzter/chatterm/types"
	"github.com/zigzter/chatterm/utils"
)

type ChatSettings struct {
	EmoteOnly     bool
	FollowersOnly bool
	SubOnly       bool
	Slow          string
}

type ChatModel struct {
	msgChan          chan types.ParsedIRCMessage
	chatContent      string
	input            string
	textinput        textinput.Model
	viewport         viewport.Model
	width            int
	height           int
	channel          string
	WsClient         *utils.WebSocketClient
	ac               *utils.Trie
	infoview         viewport.Model
	shouldRenderInfo bool
	messages         []types.ChatMessage
	username         string
	color            string
	channelUserType  string
	chatSettings     ChatSettings
}

func InitialChatModel(width int, height int) ChatModel {
	vp := viewport.New(width-2, height-5)
	vp.SetContent("")
	ip := viewport.New((width/2)-2, height-5)
	ip.SetContent("")
	utils.InitConfig()
	username := viper.GetString("username")
	color := viper.GetString("color")
	channelUserType := viper.GetString("channel-user-type")
	oauth := fmt.Sprintf("oauth:%s", viper.GetString("token"))
	channel := viper.GetString("channel")
	msgChan := make(chan types.ParsedIRCMessage, 100)
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
		input:            "",
		textinput:        ti,
		viewport:         vp,
		msgChan:          msgChan,
		WsClient:         wsClient,
		channel:          channel,
		ac:               &utils.Trie{Root: utils.NewTrieNode()},
		width:            width,
		height:           height,
		infoview:         ip,
		shouldRenderInfo: false,
		messages:         make([]types.ChatMessage, 0),
		username:         username,
		color:            color,
		channelUserType:  channelUserType,
		chatSettings: ChatSettings{
			Slow: "0",
		},
	}
}

var autocompletePrefixes = [4]string{"/ban ", "/unban ", "/info ", "@"}

// shouldAutocomplete confirms whether autocomplete should trigger,
// and returns the prefix for re-use when setting textinput value.
func shouldAutocomplete(input string) (bool, string) {
	for _, prefix := range autocompletePrefixes {
		if strings.HasPrefix(input, prefix) {
			return true, prefix
		}
	}
	return false, ""
}

// processChatInput takes in user input and determines whether the input is a command.
// If it is a command, format the command and any potential arguments
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

func listenToWebSocket(msgChan <-chan types.ParsedIRCMessage) tea.Cmd {
	return func() tea.Msg {
		return <-msgChan
	}
}

func (m *ChatModel) WrapMessages() {
	newlyWrappedChat := ""
	for _, chat := range m.messages {
		newlyWrappedChat += utils.FormatChatMessage(chat, m.viewport.Width)
	}
	m.chatContent = newlyWrappedChat
	m.viewport.SetContent(m.chatContent)
	m.viewport.GotoBottom()
}

func (m ChatModel) Init() tea.Cmd {
	return listenToWebSocket(m.msgChan)
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlX:
			m.shouldRenderInfo = false
			m.viewport.Width = m.width - 2
			m.WrapMessages()
		case tea.KeyEsc:
			m.WsClient.Conn.Close()
			return ChangeView(m, ChannelInputState)
		case tea.KeyTab:
			input := m.textinput.Value()
			if valid, prefix := shouldAutocomplete(input); valid {
				suggestion := m.ac.UpdateSuggestion(input[len(prefix):])
				m.textinput.SetValue(prefix + suggestion)
				m.textinput.CursorEnd()
			}
		case tea.KeyEnter:
			message := m.textinput.Value()
			isCommand, command, args := processChatInput(message)
			if isCommand {
				var feedback string
				res, err := twitch.SendTwitchCommand(types.TwitchCommand(command), args)
				if err != nil {
					feedback = err.Error()
				} else {
					switch resp := res.(type) {
					case *types.UserBanResp:
						data := resp.Data[0]
						if data.EndTime == nil {
							feedback = fmt.Sprintf("You banned %s from the chat.\n", args[0])
						} else {
							feedback = fmt.Sprintf("You timed out %s until %s\n", args[0], data.EndTime)
						}
					case *types.UserInfo:
						m.shouldRenderInfo = true
						m.viewport.Width = (m.width / 2) - 2
						details := resp.Details
						following := resp.Following
						followingText := ""
						if following.FollowedAt != "" {
							followingText = "Following since: " + following.FollowedAt
						}
						// TODO: instead of getting icon from messages and replacing,
						// get it from the API, along with username color
						feedback := fmt.Sprintf(
							"User: {icon}%s\nAccount created: %s\n%s\n",
							details.DisplayName,
							details.CreatedAt,
							followingText,
						)
						icon := ""
						for _, chatMsg := range m.messages {
							if chatMsg.DisplayName == args[0] {
								// TODO: move this out of the loop
								nameColor := lipgloss.NewStyle().Foreground(lipgloss.Color(chatMsg.Color))
								icon = utils.GenerateIcon(chatMsg.ChannelUserType)
								feedback += wordwrap.String(
									fmt.Sprintf("%s: %s\n", nameColor.Render(chatMsg.DisplayName), chatMsg.Message),
									m.infoview.Width,
								)
							}
						}
						feedback = strings.Replace(feedback, "{icon}", icon, 1)
						m.WrapMessages()
						m.textinput.Reset()
						m.infoview.SetContent(feedback)
						return m, listenToWebSocket(m.msgChan)
					case nil:
						// TODO: find a better way to do this?
						feedback = fmt.Sprintf("Successfully ran %s command\n", command)
					}
					m.chatContent += feedback
				}
			} else {
				nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(m.color))
				time := utils.ParseTimestamp(strconv.FormatInt(time.Now().Unix(), 10))
				m.chatContent += fmt.Sprintf("[%s]%s: %s\n", time, nameStyle.Render(m.username), message)
				m.WsClient.SendMessage([]byte("PRIVMSG #" + m.channel + " :" + message))
			}
			m.viewport.SetContent(m.chatContent)
			m.textinput.Reset()
			m.ac.Prefix = ""
			return m, listenToWebSocket(m.msgChan)
		}
	case tea.WindowSizeMsg:
		if m.shouldRenderInfo {
			m.width = (msg.Width / 2) - 2
			m.viewport.Width = m.width
			m.infoview.Width = m.width
		} else {
			m.width = msg.Width - 2
			m.viewport.Width = msg.Width - 2
			m.infoview.Width = 0
		}
		m.viewport.Height = msg.Height - 5
		m.infoview.Height = msg.Height - 5
		var vpCmd tea.Cmd
		var ipCmd tea.Cmd
		m.viewport, vpCmd = m.viewport.Update(msg)
		m.infoview, ipCmd = m.infoview.Update(msg)
		return m, tea.Batch(listenToWebSocket(m.msgChan), vpCmd, ipCmd)
	case types.ParsedIRCMessage:
		width := m.viewport.Width - 4
		switch msg := msg.Msg.(type) {
		case types.ChatMessage:
			m.messages = append(m.messages, msg)
			m.chatContent += utils.FormatChatMessage(msg, width)
			m.ac.Insert(msg.DisplayName)
		case types.SubMessage:
			m.chatContent += utils.FormatSubMessage(msg, width)
		case types.SubGiftMessage:
			m.chatContent += utils.FormatGiftSubMessage(msg, width)
		case types.MysterySubGiftMessage:
			m.chatContent += utils.FormatMysteryGiftSubMessage(msg, width)
		case types.RaidMessage:
			m.chatContent += utils.FormatRaidMessage(msg, width)
		case types.AnnouncementMessage:
			m.chatContent += utils.FormatAnnouncementMessage(msg, width)
		case types.UserListMessage:
			m.ac.Populate(msg.Users)
		case types.UserStateMessage:
			utils.SaveConfig(map[string]interface{}{
				"color":             msg.Color,
				"channel-user-type": msg.ChannelUserType,
			})
			m.channelUserType = msg.ChannelUserType
		case *types.RoomStateMessage:
			if msg.ChannelID != nil {
				utils.SaveConfig(map[string]interface{}{
					"channel-id": *msg.ChannelID,
				})
			}
			if msg.EmoteOnly != nil {
				m.chatSettings.EmoteOnly = *msg.EmoteOnly
			}
			if msg.FollowersOnly != nil {
				m.chatSettings.FollowersOnly = *msg.FollowersOnly
			}
			if msg.SubOnly != nil {
				m.chatSettings.SubOnly = *msg.SubOnly
			}
			if msg.Slow != nil {
				m.chatSettings.Slow = *msg.Slow
			}
		}
		m.viewport.SetContent(m.chatContent)
		m.viewport.GotoBottom()
		return m, listenToWebSocket(m.msgChan)

	}
	var tiCmd tea.Cmd
	var vpCmd tea.Cmd
	m.textinput, tiCmd = m.textinput.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, tiCmd, vpCmd)
	return m, tea.Batch(cmds...)
}

func iconColorizer(color string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color))
}

func (m ChatModel) View() string {
	var viewportStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#8839ef")).
		Width(m.viewport.Width)
	var infoStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#8839ef")).
		Width(m.infoview.Width)
	var b strings.Builder
	infoCloseMessage := ""
	if m.shouldRenderInfo {
		infoCloseMessage = " - [Ctrl+x] close info view"
		b.WriteString(lipgloss.JoinHorizontal(
			0,
			viewportStyle.Render(m.viewport.View()),
			infoStyle.Render(m.infoview.View())) + "\n",
		)
	} else {
		infoCloseMessage = ""
		b.WriteString(viewportStyle.Render(m.viewport.View()) + "\n")
	}
	icon := utils.GenerateIcon(m.channelUserType)
	chatSettingsString := ""
	if m.chatSettings.SubOnly {
		chatSettingsString += "[Sub Only]"
	}
	if m.chatSettings.EmoteOnly {
		chatSettingsString += "[Emote Only]"
	}
	if m.chatSettings.FollowersOnly {
		chatSettingsString += "[Followers Only]"
	}
	if m.chatSettings.Slow != "0" && m.chatSettings.Slow != "" {
		chatSettingsString += fmt.Sprintf("[Slow Mode: %ss]", m.chatSettings.Slow)
	}
	b.WriteString(chatSettingsString + icon + m.textinput.View() + "\n")
	b.WriteString(helpStyle.Render("[Esc]: return to channel selection - [Ctrl+c]: quit - [tab]: autocomplete" + infoCloseMessage))
	return b.String()
}
