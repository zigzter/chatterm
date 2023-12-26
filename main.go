package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
	"github.com/zigzter/chatterm/model"
	"github.com/zigzter/chatterm/utils"
)

func main() {
	utils.InitConfig()
	username := viper.GetString("username")
	oauth := viper.GetString("oauth")
	msgChan := make(chan model.ChatMessageWrap)
	go utils.EstablishWSConnection("GothamChess", username, oauth, msgChan)
	p := tea.NewProgram(model.InitialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
