package models

import tea "github.com/charmbracelet/bubbletea"

// Changes which model the root TUI model is displaying. Used in the Update method.
func ChangeView(model tea.Model, newView AppState) (tea.Model, tea.Cmd) {
	return model, tea.Cmd(func() tea.Msg {
		return ChangeStateMsg{NewState: newView}
	})
}
