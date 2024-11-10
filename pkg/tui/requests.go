package tui

import (
	"strings"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/network"
	tea "github.com/charmbracelet/bubbletea"
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
			return nil
		case network.JoinedLobby:
			if msg.Data == "OK" {
				m.state = inLobby
				m.clientsInLobby = append(m.clientsInLobby, Client{
					Id:   "YOU",
					Name: m.Username(),
				})
				return nil
			}

			res := strings.SplitN(msg.Data, ",", 1)

			m.clientsInLobby = append(m.clientsInLobby, Client{
				Id:   res[0],
				Name: res[1],
			})

		}
		return nil
	}

}
