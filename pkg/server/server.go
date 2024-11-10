package server

import (
	"net"

	"github.com/charmbracelet/log"
)

const (
	host = "0.0.0.0"
	port = "24816"
)

func ServerAddress() string {
	return net.JoinHostPort(host, port)
}

// Global server variable
//
// Local to this package (server directory) only
var server *Server

type Server struct {
	Lobbies      []*Lobby
	Clients      []*Client
	ln           net.Listener
	totalLobbies uint
	totalClients uint
}

func NewServer(ln net.Listener) *Server {
	s := new(Server)
	s.ln = ln
	s.totalClients = 0
	s.totalLobbies = 0
	s.Clients = make([]*Client, 0)
	s.Lobbies = make([]*Lobby, 0)

	server = s

	return s
}

func (s *Server) Start() {
	for {
		// Listen for all incoming TCP connections
		conn, err := s.ln.Accept()
		if err != nil {
			log.Error("Could not accept client connection", err)
			continue
		}

		// Concurrently run each client
		go s.newClient(conn).Start()
	}
}

func (s *Server) newClient(conn net.Conn) *Client {
	c := new(Client)
	c.Id = ClientID(s.totalClients)
	c.Conn = conn

	s.totalClients += 1
	return c
}

func (s *Server) NewLobby() *Lobby {
	l := Lobby{
		Id:       LobbyID(s.totalClients),
		Clients:  make([]*Client, 0),
		State:    WaitingForPlayers,
		commands: make(chan func()),
		finished: make(chan struct{}),
	}

	s.totalLobbies += 1
	// Go through each command one by one and run it
	//
	// Commands are requests made by client
	go func() {
		for command := range l.commands {
			command()
		}
	}()

	return &l
}

func (s *Server) FindOpenLobby() *Lobby {
	for _, lobby := range s.Lobbies {
		if lobby.State == WaitingForPlayers && len(lobby.Clients) < LOBBY_SIZE {
			return lobby
		}
	}

	// No lobbies found, so make a new lobby
	lobby := s.NewLobby()
	s.Lobbies = append(s.Lobbies, lobby)

	return lobby
}
