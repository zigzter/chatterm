package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ChannelInputModel struct {
	textinput textinput.Model
}

func InitialChannelInputModel() ChannelInputModel {
	ti := textinput.New()
	ti.Placeholder = "Enter channel name..."
	ti.Focus()
	return ChannelInputModel{
		textinput: ti,
	}
}

func (m ChannelInputModel) Init() tea.Cmd {
	return nil
}

func (m ChannelInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.textinput.Reset()
			return m, nil
		}
	}
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m ChannelInputModel) View() string {
	return m.textinput.View()
}
