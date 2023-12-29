package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/utils"
)

type ChannelInputModel struct {
	textinput textinput.Model
}

func InitialChannelInputModel() ChannelInputModel {
	utils.InitConfig()
	ti := textinput.New()
	ti.Placeholder = "a_seagull"
	ti.Focus()
	return ChannelInputModel{
		textinput: ti,
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
		case tea.KeyCtrlO:
			return m, tea.Cmd(func() tea.Msg {
				return ChangeStateMsg{NewState: ConfigState}
			})
		case tea.KeyEnter:
			if m.textinput.Value() == "exit" {
				return m, tea.Quit
			}
			viper.Set("channel", m.textinput.Value())
			if err := viper.WriteConfig(); err != nil {
				fmt.Println("Error saving config:", err)
			}
			m.textinput.Reset()
			return m, nil
		}
	}
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m ChannelInputModel) View() string {
	return fmt.Sprintf(
		"Enter channel name:\n%s\n%s",
		m.textinput.View(),
		"(Type exit or press Ctrl+c to quit. Ctrl+O for options)",
	)
}
