package server

import "github.com/Fejiberglibstein/terminal-typeracer/pkg/network"

const LOBBY_SIZE = 4

type LobbyState uint8

const (
	InGame LobbyState = iota
	WaitingForPlayers
)

type LobbyID uint

// The server side representation of a client
type Lobby struct {
	Id       LobbyID
	Clients  []*Client
	State    LobbyState
	commands chan func(network.Message)
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
