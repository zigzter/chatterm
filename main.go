package main

import (
	"log"
	"os"
	"os/signal"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zigzter/chatterm/db"
	"github.com/zigzter/chatterm/models"
	"github.com/zigzter/chatterm/utils"
)

func main() {
	configPath := utils.SetupPath()
	f, err := tea.LogToFile(configPath+"/debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	defer f.Close()
	sql := db.OpenDB(configPath)
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
