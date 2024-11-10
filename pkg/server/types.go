// Server-side networking information.
//
// Contains all the information for lobbies, clients, etc

package server

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/network"
	"github.com/charmbracelet/log"
)

const LOBBY_SIZE = 4

type ClientID uint
type LobbyID uint

type LobbyState uint8

const (
	InGame LobbyState = iota
	WaitingForPlayers
)

// The server side representation of a client
type Lobby struct {
	Id      LobbyID
	Clients []*Client
	State   LobbyState
}

// The server side representation of a Connected client
type Client struct {
	Conn  net.Conn
	Id    ClientID
	Name  string
	enc   gob.Encoder
	dec   gob.Decoder
	Lobby *Lobby
}

func (c *Client) Start() {
	c.enc = *gob.NewEncoder(c.Conn)
	c.dec = *gob.NewDecoder(c.Conn)

	for {
		// Read a message from the tcp connection
		var message network.Message
		if err := c.dec.Decode(&message); err != nil {
			// TODO better error handling, don't just kill the client maybe
			log.Error("Could not read message from client")
			break
		}
		// TODO perhaps make handling messages concurrent? idk if its necessary
		// or could cause race conditions
		c.handleMessage(message)

	}
}

func (c *Client) SendMessage(message network.Message) error {
	return c.enc.Encode(message)
}

// Used if the client was very naughty (Sending invalid requests and stuff!!)
func (c *Client) SendError(err string) {
	message := network.Message{
		Header: uint8(network.Error),
		Data:   err,
	}
	// Send the message
	if err := c.SendMessage(message); err != nil {
		log.Error("Error sending message to client", err)
	}

}

func (c *Client) handleMessage(message network.Message) {
	req := network.Request(message.Header)

	switch req {
	case network.Connect:

		// Make the length of the request only 16 chars long (client names
		// should not be over 16 chars)
		clientName := message.Data
		clientName = clientName[:16]

		server.Clients = append(server.Clients, c)
	case network.JoinLobby:
		if c.Lobby != nil {
			c.SendError("Sir you're already in a lobby")
			return
		}

		lobby := c.Lobby

		// Message to be sent to all clients already connected to the lobby
		joinedMsg := network.Message{
			Header: uint8(network.JoinedLobby),
			Data:   fmt.Sprintf("%d,%s", c.Id, c.Name),
		}

		// Let the client trying to join (c) know who is in the lobby they're
		// joining, and tell everyone in the lobby currently that someone new
		// has joined
		for _, client := range lobby.Clients {
			client.SendMessage(joinedMsg)
			// Message to be sent to the client that joined to let them know
			// what clients are connected
			joinedMsg2 := network.Message{
				Header: uint8(network.JoinedLobby),
				Data:   fmt.Sprintf("%d,%s", client.Id, client.Name),
			}
			c.SendMessage(joinedMsg2)
		}
	case network.Progress:
		// TODO
	default:
		c.SendError("Wtf man, not allowed")
	}
}
