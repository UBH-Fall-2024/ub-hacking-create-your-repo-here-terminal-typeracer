package server

import (
	"encoding/json"
	"io"
	"math/rand/v2"
	"os"
	"time"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/network"
	"github.com/charmbracelet/log"
)

const LOBBY_SIZE = 2

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

func getTypingType() string {
	jsonfile, err := os.Open("strings.json")
	if err != nil {
		return "I cannot json"
	}
	defer jsonfile.Close()

	bytes, _ := io.ReadAll(jsonfile)

	var result []string
	json.Unmarshal([]byte(bytes), &result)
	log.Print(&result)

	return result[rand.IntN(len(result)-1)]

}

func (l *Lobby) Start() {

	typer := getTypingType()

	l.State = InGame

	if err := l.SendMessage(&network.Message{
		Header: uint8(network.GameStart),
		Data:   typer,
	}); err != nil {
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
