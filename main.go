package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const cmdName = "afa-tui"

var (
	_        tea.Model = &model{}
	sockAddr string
	prompt   string
	err      string
	ver      bool
)

func init() {
	flag.StringVar(&sockAddr, "a", "", "Path to the Unix domain socket file.")
	flag.StringVar(&prompt, "p", "__AFA_PROMPT__", "Prompt string.")
	flag.StringVar(&err, "e", "__AFA_ERROR__", "Error string.")
	flag.BoolVar(&ver, "version", false, "Display version.")
	flag.Parse()
}

func main() {
	if ver {
		fmt.Printf("%s v%s (rev:%s)\n", cmdName, version, revision)
		os.Exit(0)
	}

	prog := tea.NewProgram(initialModel(sockAddr, prompt, err), tea.WithAltScreen())
	if err := prog.Start(); err != nil {
		log.Fatal(err)
	}
}
