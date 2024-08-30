package main

import (
	"fmt"
	"io"
	"net"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
)

type model struct {
	socket socket
	conn   net.Conn

	prompt bool
	status string

	loading   bool
	viewReady bool

	viewport  viewport.Model
	textinput textinput.Model
	spinner   spinner.Model

	err error
}

func initialModel() model {
	m := model{
		socket:    socket{},
		loading:   true,
		textinput: textinput.New(),
		spinner:   spinner.New(),
		prompt:    false,
		status:    "Connecting",
	}

	m.textinput.Blur()
	m.spinner.Spinner = spinner.Points
	m.spinner.Style = spinnerStyle

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.connectCmd,
		m.spinner.Tick,
	)
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
		case tea.KeyEnter:
			m.loading = true
			m.status = "Sending"
			m.textinput.Blur()
			cmds = append(cmds, m.sendCmd)
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
	case connectedMsg:
		m.conn = msg.conn
		m.prompt = true
		cmds = append(cmds, m.receiveCmd)
	case promptMsg:
		m.loading = false
		m.prompt = true
		m.textinput.Focus()
	case sentMsg:
		m.prompt = false
		m.textinput.Reset()
		m.status = "Waiting"
		cmds = append(cmds, m.receiveCmd)
	case responseMsg:
		m.loading = false
	case errMsg:
		m.err = msg
	case closeMsg:
		return m, tea.Quit
	default:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	if !m.textinput.Focused() {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Something went wrong: %s", m.err)
	}
	if !m.viewReady {
		return "Initializing..."
	}
	inputView := m.textinput.View()
	if m.loading {
		inputView = fmt.Sprintf("%s %s...", m.spinner.View(), m.status)
	}
	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.viewport.View(),
		inputView,
		m.helpView(),
	)
}

func (m model) helpView() string {
	if m.textinput.Focused() {
		return helpStyle.Render("[INSERT] enter: Submit • esc: View • ctrl+c: Quit")
	} else {
		return helpStyle.Render("[NORMAL] i/a: Prompt • j/k: Navigate • ctrl+c: Quit")
	}
}

func (m model) connectCmd() tea.Msg {
	conn, err := m.socket.dial()
	if err != nil {
		return errMsg{err}
	}
	return connectedMsg{conn}
}

func (m model) sendCmd() tea.Msg {
	input := m.textinput.Value()
	if input == "" {
		return promptMsg{}
	}
	err := m.socket.send(nil, input)
	if err != nil {
		return errMsg{err}
	}
	return sentMsg{}
}

func (m model) receiveCmd() tea.Msg {
	message, err := m.socket.receive(nil)
	if err != nil {
		if err == io.EOF {
			return closeMsg{err}
		} else {
			return errMsg{err}
		}
	}
	if m.prompt {
		return promptMsg{}
	} else {
		return responseMsg{message}
	}
	// if message == "PROMPT" {
	// 	return promptMsg{}
	// } else {
	// 	return responseMsg{message}
	// }
}
