package tui

import (
	"encoding/gob"
	"net"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/network"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type Model struct {
	width    int
	height   int
	renderer *lipgloss.Renderer
	enc      *gob.Encoder
	dec      *gob.Decoder
	Sess     ssh.Session
}

func NewModel(
	renderer *lipgloss.Renderer,
	pty *ssh.Pty,
	conn *net.Conn,
	sess ssh.Session,
) *Model {
	return &Model{
		renderer: renderer,
		width:    pty.Window.Width,
		height:   pty.Window.Height,
		enc:      gob.NewEncoder(*conn),
		dec:      gob.NewDecoder(*conn),
		Sess:     sess,
	}
}

// ReadMessage an event message from the server
func (m *Model) ReadMessage() (network.Message, error) {
	var message network.Message
	if err := m.dec.Decode(&message); err != nil {
		return network.Message{}, nil
	}
	return message, nil
}

func (m *Model) SendMessage(message *network.Message) error {
	return m.enc.Encode(message)
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

	case network.Message:
	}

	return m, nil
}

func (m *Model) View() string {
	return "HII!!!!!"
}
