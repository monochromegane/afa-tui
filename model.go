package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	viewReady bool
	viewport  viewport.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		roundedBorderSize := 2
		if !m.viewReady {
			m.viewport = viewport.New(msg.Width-roundedBorderSize, msg.Height-roundedBorderSize)
			m.viewport.Style = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				PaddingRight(2)
			m.viewport.Width = msg.Width - roundedBorderSize
			m.viewport.Height = msg.Height - roundedBorderSize
			m.viewport.SetContent("Hi, viewport")
			m.viewReady = true
		} else {
			m.viewport.Width = msg.Width - roundedBorderSize
			m.viewport.Height = msg.Height - roundedBorderSize
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.viewReady {
		return "Hi"
	}
	return fmt.Sprintf("%s", m.viewport.View())
}
