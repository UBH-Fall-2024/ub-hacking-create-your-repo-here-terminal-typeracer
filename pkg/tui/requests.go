package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/network"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type Client struct {
	Id   string
	Name string
	prog *int
}

func (m *Model) handleEvent(msg network.Message) tea.Cmd {
	return func() tea.Msg {
		switch network.Event(msg.Header) {
		case network.Error:
			m.error = &msg.Data
			return 0
		case network.JoinedLobby:
			if msg.Data == "OK" {
				m.state = inLobby
				m.clientsInLobby = append(m.clientsInLobby, &Client{
					Id:   "YOU",
					Name: m.Username(),
				})
				return 0
			}

			res := strings.SplitN(msg.Data, ",", 2)

			m.clientsInLobby = append(m.clientsInLobby, &Client{
				Id:   res[0],
				Name: res[1],
			})
			return 0

		case network.GameStart:
			m.typingInfo = &typingInfo{
				text:              msg.Data,
				correctCharacters: 0,
				typoCharacters:    0,
			}
			m.state = inGame

			return 0

		case network.LeftLobby:
			for i, client := range m.clientsInLobby {
				// Check if the client being iterated over has the same address as c
				if client.Id == msg.Data {
					// remove the client from the list
					m.clientsInLobby[i] = m.clientsInLobby[len(m.clientsInLobby)-1]
					m.clientsInLobby = m.clientsInLobby[:len(m.clientsInLobby)-1]
					break
				}
			}
			return 0
		case network.ProgUpdate:
			res := strings.SplitN(msg.Data, ",", 2)

			for _, client := range m.clientsInLobby {
				if client.Id == res[0] {
					res, err := strconv.Atoi(res[1])
					if err != nil {
						log.Print("bruh you suck")
						return nil
					}

					client.prog = &res
					return 0
				}
			}
		case network.ProgressPls:
			prog := int((float64(m.typingInfo.correctCharacters) / float64(len(m.typingInfo.text))) * 100)

			m.SendMessage(&network.Message{
				Header: uint8(network.Progress),
				Data:   fmt.Sprintf("%d", prog),
			})

			for _, client := range m.clientsInLobby {
				if client.Id == "YOU" {
					client.prog = &prog
				}
			}
		default:
			panic("unexpected network.Event")
		}
		return nil
	}

}
