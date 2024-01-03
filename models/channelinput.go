package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/utils"
)

type ChannelInputModel struct {
	textinput    textinput.Model
	authRequired bool
}

func InitialChannelInputModel() ChannelInputModel {
	authRequired := utils.InitConfig()
	ti := textinput.New()
	ti.Placeholder = "a_seagull"
	configChannel := viper.GetString("channel")
	if configChannel != "" {
		ti.SetValue(configChannel)
	}
	ti.Focus()
	return ChannelInputModel{
		textinput:    ti,
		authRequired: authRequired,
	}
}

func (m ChannelInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ChannelInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
			// TODO: re-enable config view
			// case tea.KeyCtrlO:
			// 	return ChangeView(m, ConfigState)
		case tea.KeyCtrlA:
			return ChangeView(m, AuthState)
		case tea.KeyEnter:
			if m.textinput.Value() == "exit" {
				return m, tea.Quit
			}
			utils.SaveConfig(map[string]interface{}{
				"channel": m.textinput.Value(),
			})
			return ChangeView(m, ChatState)
		}
	}
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m ChannelInputModel) View() string {
	var authMessage string
	if m.authRequired {
		authMessage = "Authentication required. Press [Ctrl+A] to start."
	}
	return fmt.Sprintf(
		"Enter channel name:\n%s\n%s\n%s",
		m.textinput.View(),
		"(Type exit or press Ctrl+c to quit. Ctrl+O for options)",
		authMessage,
	)
}
