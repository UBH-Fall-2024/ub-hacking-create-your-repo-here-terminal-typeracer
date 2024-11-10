package server

import (
	"time"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/network"
	"github.com/charmbracelet/log"
)

const LOBBY_SIZE = 4

type LobbyState uint8

const (
	WaitingForPlayers LobbyState = iota
	InGame
)

type LobbyID uint

// The server side representation of a client
type Lobby struct {
	Id       LobbyID
	Clients  []*Client
	State    LobbyState
	commands chan func()
	// Channel to send to the progress ticker when the game is over
	finished chan struct{}
}

func (l *Lobby) SendMessage(message *network.Message) error {
	for _, client := range l.Clients {
		if err := client.SendMessage(message); err != nil {
			return err
		}
	}
	return nil
}

func NewLobby() *Lobby {
	l := new(Lobby)
	l.commands = make(chan func())
	l.State = WaitingForPlayers

	l.HandleMessages()
	return l
}

// Will handle each command sent to it one by one so no race
// conditions!!!!!!!!!!!!!
func (l *Lobby) HandleMessages() {
	// Do this so we do each command individually so there's no race conditions
	for command := range l.commands {
		command()
	}
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

	go func() {
		for {
			select {
			case <-ticker.C:
				l.SendMessage(&network.Message{
					Header: uint8(network.ProgressPls),
					Data:   "",
				})
			case <-l.finished:
				ticker.Stop()
				return
			}
		}
	}()

}
