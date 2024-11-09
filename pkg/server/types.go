// Server-side networking information.
//
// Contains all the information for lobbies, clients, etc

package server

import (
	"encoding/gob"
	"net"
)

type ClientID uint
type LobbyID uint

type Lobby struct {
	Id LobbyID
}

type Client struct {
	Conn net.Conn
	Id   ClientID
	enc  gob.Encoder
	dec  gob.Decoder
}

func (c *Client) Start() {
	c.enc = *gob.NewEncoder(c.Conn)
	c.dec = *gob.NewDecoder(c.Conn)

	for {

	}
}
