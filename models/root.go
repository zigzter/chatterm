package models

import tea "github.com/charmbracelet/bubbletea"

type AppState int

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
	return RootModel{
		State: 0,
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.State {
	case ChannelInputState:
		return m.ChannelInput.Update(msg)
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
