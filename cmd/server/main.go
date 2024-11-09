package main

import (
	"github.com/charmbracelet/log"
	"net"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/server"
)

const (
	host = "0.0.0.0"
	port = "24816"
)

func main() {
	ln, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		log.Fatal(err)
	}

	server.NewServer(ln).Start()

}
