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

type Lobby struct {
	Id LobbyID
}

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
		var message network.Message
		if err := c.dec.Decode(&message); err != nil {
			log.Error("Could not decode")
		}
		go c.handleMessage(message)

	}
}

func (c *Client) handleMessage(message network.Message) {

}
