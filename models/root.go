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
	case ChangeStateMsg:
		m.State = msg.NewState
		// This is to handle pressing enter on a pre-filled ChannelInput value
		// TODO: find a better way to do this
		if m.State == ChatState && !m.IsChatInitialized {
			m.Chat = InitialChatModel()
			chatCmd := m.Chat.Init()
			m.IsChatInitialized = true
			return m, chatCmd
		}
		return m, nil
	}
	switch m.State {
	case ChannelInputState:
		newModel, cmd := m.ChannelInput.Update(msg)
		m.ChannelInput = newModel.(ChannelInputModel)
		return m, cmd
	case ChatState:
		if !m.IsChatInitialized {
			m.Chat = InitialChatModel()
			initCmd := m.Chat.Init()
			m.IsChatInitialized = true
			return m, initCmd
		}
		newModel, cmd := m.Chat.Update(msg)
		m.Chat = newModel.(ChatModel)
		return m, cmd
	case ConfigState:
		newModel, cmd := m.Config.Update(msg)
		m.Config = newModel.(ConfigModel)
		return m, cmd
	case AuthState:
		newModel, cmd := m.Auth.Update(msg)
		m.Auth = newModel.(AuthModel)
		return m, cmd
	default:
		return m, nil
	}
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
