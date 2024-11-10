package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

func (m *Model) RenderClients() string {

	clients := make([]string, 0)

	for _, client := range m.clientsInLobby {
		var name string
		if client.Id == "YOU" {
			name = client.Name + " (you)"
		} else {
			name = client.Name
		}

		clients = append(clients, name)

	}

	l := list.New(clients).
		Enumerator(list.Bullet).
		EnumeratorStyle(lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("93")).
			MarginRight(1)).ItemStyle(lipgloss.NewStyle().Width(22))

	return l.String()
}


func (m* Model) RenderTyper() string {

}
