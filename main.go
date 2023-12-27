package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zigzter/chatterm/model"
)

func main() {
	p := tea.NewProgram(model.InitialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
