package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	viewReady bool

	viewport  viewport.Model
	textinput textinput.Model
}

func initialModel() model {
	m := model{
		textinput: textinput.New(),
	}

	m.textinput.Focus()

	return m
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
		textinputHeight := 1
		helpHeight := 1
		footerHeight := textinputHeight + helpHeight + roundedBorderSize

		if !m.viewReady {
			m.viewport = viewport.New(msg.Width-roundedBorderSize, msg.Height-footerHeight)
			m.viewport.Style = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62")).
				PaddingRight(2)
			m.viewport.Width = msg.Width - roundedBorderSize
			m.viewport.Height = msg.Height - footerHeight
			m.viewport.SetContent("Hi, viewport")
			m.viewReady = true
		} else {
			m.viewport.Width = msg.Width - roundedBorderSize
			m.viewport.Height = msg.Height - footerHeight
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.viewReady {
		return "Hi"
	}
	return fmt.Sprintf(
		"%s\n%s",
		m.viewport.View(),
		m.textinput.View(),
	)
}
