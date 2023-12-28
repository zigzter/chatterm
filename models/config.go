package models

import tea "github.com/charmbracelet/bubbletea"

type ConfigModel struct{}

func InitialConfigModel() ConfigModel {
	return ConfigModel{}
}

func (m ConfigModel) Init() tea.Cmd {
	return nil
}

func (m ConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m ConfigModel) View() string {
	return ""
}
