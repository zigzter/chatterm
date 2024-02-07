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
	SettingsState
	ChatState
	AuthState
)

type RootModel struct {
	State             AppState
	Chat              ChatModel
	ChannelInput      ChannelInputModel
	Settings          SettingsModel
	Auth              AuthModel
	IsChatInitialized bool
	Width             int
	Height            int
}

func InitialRootModel() RootModel {
	channelInputModel := InitialChannelInputModel()
	settingsModel := InitialSettingsModel()
	authModel := InitialAuthModel()
	return RootModel{
		State:             0,
		ChannelInput:      channelInputModel,
		Settings:          settingsModel,
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
		m.Width = msg.Width
		m.Height = msg.Height
	case ChangeStateMsg:
		m.State = msg.NewState
		// This is to handle pressing enter on a pre-filled ChannelInput value
		// TODO: find a better way to do this
		if m.State == ChatState {
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
	case SettingsState:
		newModel, newCmd := m.Settings.Update(msg)
		m.Settings = newModel.(SettingsModel)
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
	case SettingsState:
		return m.Settings.View()
	case AuthState:
		return m.Auth.View()
	default:
		return ""
	}
}
