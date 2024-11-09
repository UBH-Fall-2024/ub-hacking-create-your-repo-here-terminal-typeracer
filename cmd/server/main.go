package main

import (
	"log"
	"net"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/network"
)

const (
	host = "0.0.0.0"
	port = "24816"
)

type Server struct {
	Lobbies      []network.Lobby  //TODO
	Clients      []network.Client //TODO
	ln           net.Listener
	totalLobbies uint
	totalClients uint
}

func main() {
	ln, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Fatal(err)
	}

	server := Server{
		Lobbies:      []network.Lobby{},  //TODO
		Clients:      []network.Client{}, //TODO
		ln:           ln,
		totalLobbies: 0,
		totalClients: 0,
	}

	startServer(server)

}

func startServer(server Server) {
	for {
		conn, err := server.ln.Accept()
		if err != nil {
			log.Print("Could not accept client connection: ", err)
			continue
		}

		go startClient(conn)
	}
}

func startClient(conn net.Conn) {

}

func (s *Server) newClient(conn net.Conn) network.Client {
	currentClients := s.totalClients

	s.totalClients += 1

	return network.Client{
		Id:   network.ClientID(currentClients),
		Conn: conn,
	}
}
