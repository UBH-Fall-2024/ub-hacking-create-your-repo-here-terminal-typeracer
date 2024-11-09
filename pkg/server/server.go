package server

import (
	"net"

	"github.com/charmbracelet/log"
)

type Server struct {
	Lobbies      []Lobby  //TODO
	Clients      []Client //TODO
	ln           net.Listener
	totalLobbies uint
	totalClients uint
}

func NewServer(ln net.Listener) *Server {
	s := new(Server)
	s.ln = ln
	s.totalClients = 0
	s.totalLobbies = 0
	s.Clients = make([]Client, 8)
	s.Lobbies = make([]Lobby, 8)

	return s
}

func (s *Server) StartServer() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Error("Could not accept client connection", err)
			continue
		}

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
