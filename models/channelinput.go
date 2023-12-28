package models

import tea "github.com/charmbracelet/bubbletea"

type ChannelInputModel struct{}

func InitialChannelInputModel() ChannelInputModel {
	return ChannelInputModel{}
}

func (m ChannelInputModel) Init() tea.Cmd {
	return nil
}

func (m ChannelInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ChannelInputModel) View() string {
	return ""
}
