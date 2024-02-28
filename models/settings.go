package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/utils"
)

type SettingsModel struct {
	form *huh.Form
}

// TODO: Implement logic
func isHexColor(input string) error {
	return nil
}

var (
	showBadges        bool   = true
	showTimestamps    bool   = true
	highlightSubs     bool   = true
	highlightRaids    bool   = true
	firstChatterColor string = "e64553"
)

type SettingsSaved struct{}

func SaveSettings(shouldSave bool) {
	if shouldSave {
		utils.SaveConfig(map[string]interface{}{
			utils.ShowBadgesKey:        showBadges,
			utils.ShowTimestampsKey:    showTimestamps,
			utils.HighlightSubsKey:     highlightSubs,
			utils.HighlightRaidsKey:    highlightRaids,
			utils.FirstChatterColorKey: firstChatterColor,
		})
	}
}

func InitialSettingsModel() SettingsModel {
	m := SettingsModel{}
	showBadges = viper.GetBool(utils.ShowBadgesKey)
	showTimestamps = viper.GetBool(utils.ShowTimestampsKey)
	highlightSubs = viper.GetBool(utils.HighlightSubsKey)
	highlightRaids = viper.GetBool(utils.HighlightRaidsKey)
	firstChatterColor = viper.GetString(utils.FirstChatterColorKey)
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[bool]().
				Key("timestamps").
				Options(
					huh.NewOption("Show", true),
					huh.NewOption("Hide", false),
				).
				Title("Display timestamps").
				Description("Choose whether timestamps should be shown").
				Value(&showTimestamps),
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
					SaveSettings(v)
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
	if m.form.State == huh.StateCompleted {
		return ChangeView(m, ChannelInputState)
	}
	return m, tea.Batch(cmds...)
}

func (m SettingsModel) View() string {
	return m.form.View()
}
