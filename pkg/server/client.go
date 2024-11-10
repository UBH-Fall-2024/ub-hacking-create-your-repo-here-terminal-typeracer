// Server-side networking information.
//
// Contains all the information for lobbies, clients, etc

package server

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"

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
	commands chan func(network.Message)
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
			c.Disconnect()
			break
		}
		// TODO perhaps make handling messages concurrent? idk if its necessary
		// or could cause race conditions
		c.handleMessage(message)

	}
}

func (c *Client) SendMessage(message *network.Message) error {
	return c.enc.Encode(message)
}

// Used if the client was very naughty (Sending invalid requests and stuff!!)
func (c *Client) SendError(err string) error {
	message := network.Message{
		Header: uint8(network.Error),
		Data:   err,
	}
	return c.SendMessage(&message)
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
			client.SendMessage(&joinedMsg)
			// Message to be sent to the client that joined to let them know
			// what clients are connected
			joinedMsg2 := network.Message{
				Header: uint8(network.JoinedLobby),
				Data:   fmt.Sprintf("%d,%s", client.Id, client.Name),
			}
			c.SendMessage(&joinedMsg2)
		}

		if len(lobby.Clients) >= LOBBY_SIZE {
			lobby.Start()
		}
	case network.Progress:
		// TODO
	default:
		c.SendError("Wtf man, not allowed")
	}
}

func (c *Client) Disconnect() {

	lobby := c.Lobby
	c.Lobby = nil

	if err := c.SendError("Disconnected from server"); err != nil {
		log.Print("Could not tell the client they suck")
	}

	for i, client := range lobby.Clients {
		// Check if the client being iterated over has the same address as c
		if client == c {
			// remove the client from the list
			lobby.Clients[i] = lobby.Clients[len(lobby.Clients)-1]
			lobby.Clients = lobby.Clients[:len(lobby.Clients)-1]
			break
		}
	}

	message := network.Message{
		Header: uint8(network.LeftLobby),
		Data:   string(c.Id),
	}

	// Doing double for loop because i am scared removing from list will mess
	// with iteration
	for _, client := range lobby.Clients {
		client.SendMessage(&message)
	}
}

func (l *Lobby) SendMessage(message *network.Message) error {
	for _, client := range l.Clients {
		if err := client.SendMessage(message); err != nil {
			return err
		}
	}
	return nil
}

func (l *Lobby) Start() {

	l.State = InGame

	m := network.Message{
		Header: uint8(network.GameStart),
		Data:   "type this right now or else you will lose the entire game and i will cry and hate you forever because you suck at typing and you really need to do better and eat some chocolate chocolate chip ice cream cone tasty mm delicious",
	}

	if err := l.SendMessage(&m); err != nil {
		log.Error("Could not start the server")
	}

	ticker := time.NewTicker(1 * time.Second)

	quit := make(chan struct {})

	go func() {

	}

}