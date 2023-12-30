package main

import (
	"log"
	"os"
	"os/signal"

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
	p := tea.NewProgram(m)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	go func() {
		<-shutdown
		m.Chat.WsClient.Conn.Close()
		os.Exit(0)
	}()
}
