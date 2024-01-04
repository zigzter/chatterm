package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

type AppState int

type ChangeStateMsg struct {
	NewState AppState
}

const (
	ChannelInputState AppState = iota
	ConfigState
	ChatState
	AuthState
)

type RootModel struct {
	State             AppState
	Chat              ChatModel
	ChannelInput      ChannelInputModel
	Config            ConfigModel
	Auth              AuthModel
	IsChatInitialized bool
	Width             int
	Height            int
}

func InitialRootModel() RootModel {
	channelInputModel := InitialChannelInputModel()
	configModel := InitialConfigModel()
	authModel := InitialAuthModel()
	return RootModel{
		State:             0,
		ChannelInput:      channelInputModel,
		Config:            configModel,
		Auth:              authModel,
		IsChatInitialized: false,
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width, m.Height = msg.Width, msg.Height
		return m, nil
	case ChangeStateMsg:
		m.State = msg.NewState
		// This is to handle pressing enter on a pre-filled ChannelInput value
		// TODO: find a better way to do this
		if m.State == ChatState && !m.IsChatInitialized {
			m.Chat = InitialChatModel(m.Width, m.Height)
			chatCmd := m.Chat.Init()
			m.IsChatInitialized = true
			return m, chatCmd
		}
		return m, nil
	}
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch m.State {
	case ChannelInputState:
		newModel, newCmd := m.ChannelInput.Update(msg)
		m.ChannelInput = newModel.(ChannelInputModel)
		cmd = newCmd
	case ChatState:
		newModel, newCmd := m.Chat.Update(msg)
		m.Chat = newModel.(ChatModel)
		cmd = newCmd
	case ConfigState:
		newModel, newCmd := m.Config.Update(msg)
		m.Config = newModel.(ConfigModel)
		cmd = newCmd
	case AuthState:
		newModel, newCmd := m.Auth.Update(msg)
		m.Auth = newModel.(AuthModel)
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	switch m.State {
	case ChannelInputState:
		return m.ChannelInput.View()
	case ChatState:
		return m.Chat.View()
	case ConfigState:
		return m.Config.View()
	case AuthState:
		return m.Auth.View()
	default:
		return ""
	}
}
