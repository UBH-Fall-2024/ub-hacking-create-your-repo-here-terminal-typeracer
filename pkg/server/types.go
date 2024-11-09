// Server-side networking information.
//
// Contains all the information for lobbies, clients, etc

package server

import (
	"encoding/gob"
	"net"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/network"
	"github.com/charmbracelet/log"
)

type ClientID uint
type LobbyID uint

// The server side representation of a client
type Lobby struct {
	Id      LobbyID
	Clients []Client
}

// The server side representation of a Connected client
type Client struct {
	Conn  net.Conn
	Id    ClientID
	enc   gob.Encoder
	dec   gob.Decoder
	lobby *Lobby
}

func (c *Client) Start() {
	c.enc = *gob.NewEncoder(c.Conn)
	c.dec = *gob.NewDecoder(c.Conn)

	for {
		// Read a message from the tcp connection
		var message network.Message
		if err := c.dec.Decode(&message); err != nil {
			log.Error("Could not decode")
			continue
		}
		// TODO perhaps make handling messages concurrent? idk if its necessary
		// or could cause race conditions
		c.handleMessage(message)

	}
}

func (c *Client) handleMessage(message network.Message) {
}
