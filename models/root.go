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
)

type RootModel struct {
	State        AppState
	Chat         ChatModel
	ChannelInput ChannelInputModel
	Config       ConfigModel
}

func InitialRootModel() RootModel {
	channelInputModel := InitialChannelInputModel()
	chatModel := InitialChatModel()
	configModel := InitialConfigModel()
	return RootModel{
		State:        0,
		ChannelInput: channelInputModel,
		Chat:         chatModel,
		Config:       configModel,
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ChangeStateMsg:
		m.State = msg.NewState
		return m, nil
	}
	switch m.State {
	case ChannelInputState:
		newModel, cmd := m.ChannelInput.Update(msg)
		m.ChannelInput = newModel.(ChannelInputModel)
		return m, cmd
	case ChatState:
		return m.Chat.Update(msg)
	case ConfigState:
		return m.Config.Update(msg)
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
	default:
		return ""
	}
}
