package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

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
			content := ""
			for i := range 100 {
				content += fmt.Sprintf(" %03d: Hi, viewport\n", i)
			}
			m.viewport.SetContent(content)
			m.viewReady = true
		} else {
			m.viewport.Width = msg.Width - roundedBorderSize
			m.viewport.Height = msg.Height - footerHeight
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			m.textinput.Blur()
		default:
			if m.textinput.Focused() {
				m.textinput, cmd = m.textinput.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				if key := msg.String(); key == "i" || key == "a" {
					m.textinput.Focus()
				}
			}
		}
	}

	if !m.textinput.Focused() {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.viewReady {
		return "Initilizing..."
	}
	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.viewport.View(),
		m.textinput.View(),
		m.helpView(),
	)
}

func (m model) helpView() string {
	if m.textinput.Focused() {
		return helpStyle("[INSERT] enter: Submit • esc: View • ctrl+c: Quit")
	} else {
		return helpStyle("[NORMAL] i/a: Edit • j/k: Navigate • ctrl+c: Quit")
	}
}
