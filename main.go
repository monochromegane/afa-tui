package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &model{}

func main() {
	prog := tea.NewProgram(model{}, tea.WithAltScreen())
	if err := prog.Start(); err != nil {
		log.Fatal(err)
	}
}
