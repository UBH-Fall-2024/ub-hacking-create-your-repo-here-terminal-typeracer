package tui

import (
	"encoding/gob"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type Model struct {
	renderer *lipgloss.Renderer
	width    int
	height   int
	enc      *gob.Encoder
	dec      *gob.Decoder
}

func NewModel(renderer *lipgloss.Renderer, pty *ssh.Pty, conn *net.Conn) *Model {
	return &Model{
		renderer: renderer,
		width:    pty.Window.Width,
		height:   pty.Window.Height,
		enc:      gob.NewEncoder(*conn),
		dec:      gob.NewDecoder(*conn),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

	}

	return m, nil
}

func (m *Model) View() string {
	return "HII!!!!!"
}
