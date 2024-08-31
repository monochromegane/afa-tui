package main

import (
	"flag"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	_        tea.Model = &model{}
	sockAddr string
	prompt   string
)

func init() {
	flag.StringVar(&sockAddr, "a", "", "Path to the Unix domain socket file.")
	flag.StringVar(&prompt, "p", "__AFA_PROMPT__", "Prompt string.")
	flag.Parse()
}

func main() {
	prog := tea.NewProgram(initialModel(sockAddr, prompt), tea.WithAltScreen())
	if err := prog.Start(); err != nil {
		log.Fatal(err)
	}
}
