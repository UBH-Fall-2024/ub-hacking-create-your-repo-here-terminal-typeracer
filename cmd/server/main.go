package main

import (
	"log"
	"net"
)

const (
	host = "0.0.0.0"
	port = "24816"
)

type Server struct {
	lobbies []Lobby  //TODO
	clients []Client //TODO
	ln      net.Listener
}

func main() {
	ln, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Fatal(err)
	}

	server := Server{
		lobbies: []invalid{}, //TODO
		clients: []invalid{}, //TODO
		ln:      ln,
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
