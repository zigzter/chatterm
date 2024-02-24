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

var (
	showBadges        bool
	highlightSubs     bool
	highlightRaids    bool
	firstChatterColor string
)

func InitialSettingsModel() SettingsModel {
	m := SettingsModel{}
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Key("badges").
				Options(
					huh.NewOption("Show", true),
					huh.NewOption("Hide", false),
				).
				Title("Display badges").
				Description("Choose whether mod/vip badges should be shown").
				Value(&showBadges),
			huh.NewSelect[bool]().
				Key("subs").
				Options(
					huh.NewOption("Show", true),
					huh.NewOption("Hide", false),
				).
				Title("Highlight subscriptions").
				Description("Choose whether to highlight subscription messages").
				Value(&highlightSubs),
			huh.NewSelect[bool]().
				Key("raids").
				Options(
					huh.NewOption("Show", true),
					huh.NewOption("Hide", false),
				).
				Title("Highlight raids").
				Description("Choose whether to highlight raid messages").
				Value(&highlightRaids),
			huh.NewInput().
				Title("First chatter color").
				Prompt("#").
				Validate(isHexColor).
				Value(&firstChatterColor),
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
		case "esc":
			return ChangeView(m, ChannelInputState)
		case "ctrl+c":
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
