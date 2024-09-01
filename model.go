package main

import (
	"encoding/gob"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	StateConnecting = "Connecting"
	StateSending    = "Sending"
	StateReceiving  = "Receiving"

	MessageInitial = "How can I assist you today?"

	HelpNormalMode         = "  [NORMAL] i/a: Prompt • j/k: Navigate • ctrl+c: Quit"
	HelpNormalReadOnlyMode = "  [NORMAL](Read Only) j/k: Navigate • ctrl+c: Quit"
	HelpInsertMode         = "  [INSERT] enter: Submit • esc: View • ctrl+c: Quit"
)

var (
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
)

type model struct {
	socket  socket
	encoder *gob.Encoder
	decoder *gob.Decoder

	prompt string
	state  string

	loading     bool
	interacting bool
	viewReady   bool
	chunk       int

	viewport  viewport.Model
	content   content
	textinput textinput.Model
	spinner   spinner.Model

	err error
}

func initialModel(sockAddr, prompt string) model {
	m := model{
		socket:      socket{sockAddr},
		loading:     true,
		interacting: false,
		textinput:   textinput.New(),
		spinner:     spinner.New(),
		prompt:      prompt,
		state:       StateConnecting,
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
	var err error

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		roundedBorderSize := 2
		paddingSize := 2
		textinputHeight := 1
		helpHeight := 1
		footerHeight := textinputHeight + helpHeight + roundedBorderSize

		if !m.viewReady {
			m.viewport = viewport.New(msg.Width-roundedBorderSize, msg.Height-footerHeight)
			m.viewport.Style = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62"))
			m.viewport.SetContent(MessageInitial)
			m.viewport.Width = msg.Width - roundedBorderSize
			m.viewport.Height = msg.Height - footerHeight
			m.content, err = newContent(
				m.viewport.Width-roundedBorderSize-paddingSize,
				m.viewport.Width-roundedBorderSize,
			)
			m.viewReady = true
		} else {
			m.viewport.Width = msg.Width - roundedBorderSize
			m.viewport.Height = msg.Height - footerHeight
			m.content, err = m.content.update(
				m.viewport.Width-roundedBorderSize-paddingSize,
				m.viewport.Width-roundedBorderSize,
			)
			m.viewport.SetContent(m.content.rendered)
		}
		if err != nil {
			m.err = err
			return m, m.errCmd
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			m.textinput.Blur()
		case tea.KeyEnter:
			message := m.textinput.Value()
			if message == "" || m.loading {
				break
			}
			m.content, err = m.content.add(fmt.Sprintf("# You\n\n%s\n", message))
			if err != nil {
				m.err = err
				return m, m.errCmd
			}
			m.viewport.SetContent(m.content.rendered)

			m.loading = true
			m.textinput.Blur()
			m.state = StateSending
			cmds = append(cmds, m.sendCmd)
		default:
			if m.textinput.Focused() {
				m.textinput, cmd = m.textinput.Update(msg)
				cmds = append(cmds, cmd)
			} else {
				if key := msg.String(); m.interacting && (key == "i" || key == "a") {
					m.textinput.Focus()
				}
			}
		}
	case connectedMsg:
		m.encoder = gob.NewEncoder(msg.conn)
		m.decoder = gob.NewDecoder(msg.conn)
		cmds = append(cmds, m.receiveCmd)
	case promptMsg:
		m.interacting = true
		m.loading = false
		m.chunk = 0
	case sentMsg:
		m.textinput.Reset()
		m.state = StateReceiving
		cmds = append(cmds, m.receiveCmd)
	case responseMsg:
		m.chunk += 1
		format := "%s"
		if m.interacting && m.chunk == 1 {
			format = "# Assistant\n\n%s"
		}
		m.content, err = m.content.add(fmt.Sprintf(format, msg.message))
		if err != nil {
			m.err = err
			return m, m.errCmd
		}
		m.viewport.SetContent(m.content.rendered)
		m.viewport.GotoBottom()

		cmds = append(cmds, m.receiveCmd)
	case errMsg:
		m.err = msg
	case closeMsg:
		m.loading = false
		m.interacting = false
		m.chunk = 0
		m.textinput.Blur()
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
		inputView = fmt.Sprintf("%s %s...", m.spinner.View(), m.state)
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
		return helpStyle.Render(HelpInsertMode)
	} else if m.interacting {
		return helpStyle.Render(HelpNormalMode)
	} else {
		return helpStyle.Render(HelpNormalReadOnlyMode)
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
	err := m.socket.send(m.encoder, input)
	if err != nil {
		return errMsg{err}
	}
	return sentMsg{}
}

func (m model) receiveCmd() tea.Msg {
	message, err := m.socket.receive(m.decoder)
	if err != nil {
		if err == io.EOF {
			return closeMsg{err}
		} else {
			return errMsg{err}
		}
	}
	if message == m.prompt {
		return promptMsg{}
	} else {
		return responseMsg{message}
	}
}

func (m model) errCmd() tea.Msg {
	return errMsg{m.err}
}
