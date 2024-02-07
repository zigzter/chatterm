package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type SettingsModel struct {
	form    *huh.Form
	confirm bool
}

func isHexColor(input string) error {
	return nil
}

func InitialSettingsModel() SettingsModel {
	m := SettingsModel{}
	var color string
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Key("badges").
				Options(huh.NewOptions("Show", "Hide")...).
				Title("Display badges").
				Description("Choose whether mod/vip badges should be shown"),
			huh.NewSelect[string]().
				Key("subs").
				Options(huh.NewOptions("Show", "Hide")...).
				Title("Highlight subscriptions").
				Description("Choose whether to highlight subscription messages"),
			huh.NewSelect[string]().
				Key("raids").
				Options(huh.NewOptions("Show", "Hide")...).
				Title("Highlight raids").
				Description("Choose whether to highlight raid messages"),
			huh.NewInput().
				Title("First chatter color").
				Prompt("#").
				Validate(isHexColor).
				Value(&color),
			huh.NewConfirm().
				Key("save").
				Title("Save Settings").
				Validate(func(v bool) error {
					if !v {
						m.confirm = false
						return nil
					}
					m.confirm = true
					return nil
				}).
				Affirmative("Save").
				Negative("Cancel"),
		),
	)
	return m
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	var cmds []tea.Cmd
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m SettingsModel) View() string {
	return m.form.View() + fmt.Sprintf("\nConfirm: %t", m.confirm)
}
