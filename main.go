package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zigzter/chatterm/models"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %w", err)
	}
	defer f.Close()
	m := models.InitialRootModel()
	defer m.Chat.WsClient.Conn.Close()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
