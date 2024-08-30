package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &model{}

func main() {
	addr := os.Args[1]
	prog := tea.NewProgram(initialModel(addr, "PROMPT"), tea.WithAltScreen())
	if err := prog.Start(); err != nil {
		log.Fatal(err)
	}
}
