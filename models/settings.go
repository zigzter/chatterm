package models

import (
	"errors"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/utils"
)

type SettingsModel struct {
	form *huh.Form
}

var hexRegex *regexp.Regexp = regexp.MustCompile("(?i)^#[0-9A-F]{6}$")

func isHexColor(input string) error {
	if hexRegex.MatchString(input) {
		return nil
	}
	return errors.New("Color must be a valid hex code")
}

type SettingsSaved struct{}

func InitialSettingsModel() SettingsModel {
	m := SettingsModel{}
	showBadges := viper.GetBool(utils.ShowBadgesKey)
	showTimestamps := viper.GetBool(utils.ShowTimestampsKey)
	highlightSubs := viper.GetBool(utils.HighlightSubsKey)
	highlightRaids := viper.GetBool(utils.HighlightRaidsKey)
	firstChatterColor := viper.GetString(utils.FirstTimeChatterColorKey)
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
				Key("color").
				Title("First chatter color").
				Prompt(">").
				Validate(isHexColor).
				Value(&firstChatterColor),
			huh.NewConfirm().
				Key("save").
				Title("Save Settings").
				Validate(func(v bool) error {
					m.SaveSettings(v)
					return nil
				}).
				Affirmative("Save").
				Negative("Cancel"),
		),
	)
	return m
}

func (m SettingsModel) SaveSettings(shouldSave bool) {
	if shouldSave {
		utils.SaveConfig(map[string]interface{}{
			utils.ShowBadgesKey:            m.form.GetBool("badges"),
			utils.ShowTimestampsKey:        m.form.GetBool("timestamps"),
			utils.HighlightSubsKey:         m.form.GetBool("subs"),
			utils.HighlightRaidsKey:        m.form.GetBool("raids"),
			utils.FirstTimeChatterColorKey: m.form.GetString("color"),
		})
	}
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
