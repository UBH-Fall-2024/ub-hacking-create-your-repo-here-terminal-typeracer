package server

import (
	"net"

	"github.com/charmbracelet/log"
)

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
	s.Clients = make([]*Client, 8)
	s.Lobbies = make([]*Lobby, 8)

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

func (s *Server) FindOpenLobby() *Lobby {
	for _, lobby := range s.Lobbies {
		if lobby.State == WaitingForPlayers && len(lobby.Clients) < LOBBY_SIZE {
			return lobby
		}
	}
	// No lobbies found, so make a new lobby

	lobby := new(Lobby)
	s.Lobbies = append(s.Lobbies, lobby)

	return lobby
}
