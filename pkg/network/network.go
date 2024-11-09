package network

import "net"

type ClientID uint
type LobbyID uint

type Lobby struct {
	Id LobbyID
}

type Client struct {
	Conn net.Conn
	Id   ClientID
}
