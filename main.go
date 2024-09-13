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
	err      string
)

func init() {
	flag.StringVar(&sockAddr, "a", "", "Path to the Unix domain socket file.")
	flag.StringVar(&prompt, "p", "__AFA_PROMPT__", "Prompt string.")
	flag.StringVar(&err, "e", "__AFA_ERROR__", "Error string.")
	flag.Parse()
}

func main() {
	prog := tea.NewProgram(initialModel(sockAddr, prompt, err), tea.WithAltScreen())
	if err := prog.Start(); err != nil {
		log.Fatal(err)
	}
}
