package tui

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/network"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	zone "github.com/lrstanley/bubblezone"
)

type modelState uint8

const (
	// Boring and doesn't want to play with anyone (Before connecting to any
	// lobbies)
	noConnection modelState = iota

	// While they are waiting to connect to a lobby
	connecting

	// Waiting in lobby
	inLobby

	// In game and having fun with all their friends
	inGame
)

const (
	joinLobby = "join_lobby"
)

type Model struct {
	width    int
	height   int
	renderer *lipgloss.Renderer
	enc      *gob.Encoder
	dec      *gob.Decoder
	sess     ssh.Session
	state    modelState

	// Various ui menu status
	error          *string
	clientsInLobby []Client

	zone  zone.Manager
	style Style
}

func NewModel(
	renderer *lipgloss.Renderer,
	pty *ssh.Pty,
	conn *net.Conn,
	sess ssh.Session,
) *Model {
	return &Model{
		width:          pty.Window.Width,
		height:         pty.Window.Height,
		renderer:       renderer,
		enc:            gob.NewEncoder(*conn),
		dec:            gob.NewDecoder(*conn),
		sess:           sess,
		state:          noConnection,
		error:          nil,
		clientsInLobby: make([]Client, 0),
		zone:           *zone.New(),
		style:          NewStyle(renderer),
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

func (m *Model) Username() string {
	return m.sess.User()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			if msg.String() == "q" && m.error != nil {
				cmds = append(cmds, func() tea.Msg {
					m.error = nil
					return nil
				})
			}
		}
	case tea.MouseMsg:
		// No hover...
		if msg.Action != tea.MouseActionRelease || msg.Button != tea.MouseButtonLeft {
			return m, nil
		}
		log.Print(m.zone.Get(joinLobby))

		if m.zone.Get(joinLobby).InBounds(msg) {
			cmds = append(cmds, func() tea.Msg {
				m.SendMessage(&network.Message{
					Header: uint8(network.JoinLobby),
					Data:   "",
				})
				m.state = connecting
				return nil
			})

		}
	case network.Message:
		cmds = append(cmds, m.handleEvent(msg))
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var view string

	switch m.state {
	case noConnection:
		view = m.renderer.Place(
			m.width,
			m.height,
			lipgloss.Center, // h align
			lipgloss.Center, // v align
			m.zone.Mark(joinLobby, m.style.buttonStyle.Render("Connect to lobby")),
		)
	case connecting:
		view = m.renderer.Place(
			m.width,
			m.height,
			lipgloss.Center, // h align
			lipgloss.Center, // v align
			m.style.connectingStyle.Render("Connecting..."),
		)
	case inLobby:
		view = m.renderer.Place(
			m.width,
			m.height,
			lipgloss.Center, // h align
			lipgloss.Center, // v align
			m.RenderClients(),
		)
	}

	if m.error != nil {
		view = m.renderer.Place(
			m.width,
			m.height,
			lipgloss.Center, // h align
			lipgloss.Center, // v align
			m.style.errorStyle.Render(*m.error),
		)
	}
	return m.zone.Scan(view)
}
