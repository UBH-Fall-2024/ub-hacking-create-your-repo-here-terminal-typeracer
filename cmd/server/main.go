package main

import (
	"github.com/charmbracelet/log"
	"net"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/server"
)

func main() {
	ln, err := net.Listen("tcp", server.ServerAddress())
	if err != nil {
		log.Fatal(err)
	}

	server.NewServer(ln).Start()

}
