package main

import (
	"log"
	"os"
	"os/signal"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zigzter/chatterm/db"
	"github.com/zigzter/chatterm/models"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %w", err)
	}
	defer f.Close()
	sql := db.OpenDB()
	defer sql.Close()
	db.CreateTables(sql)

	m := models.InitialRootModel()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

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
