package tui

import (
	"fmt"
	"strings"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/server"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
)

var colors = []string{"202", "207", "63", "84", "226"}

func (c *Client) getName() string {
	name := c.Name
	if c.Id == "YOU" {
		name = c.Name + " (you)"
	}
	return name
}

func (m *Model) RenderClients() string {

	clients := make([]string, 0)

	for _, client := range m.clientsInLobby {
		clients = append(clients, client.getName())

	}

	l := list.New(clients).
		Enumerator(list.Bullet).
		EnumeratorStyle(lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("93")).
			MarginRight(1)).ItemStyle(lipgloss.NewStyle().Width(22)).
		String()

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.style.notiStyle.Render(fmt.Sprint("Waiting for ", server.LOBBY_SIZE-len(clients), " more player(s) to join")),
		l,
	)
}

func renderCar(color string, prog int) string {
	const length float64 = 32

	f := float64(prog) / 100.0

	start := strings.Repeat("▬", int(f*length))
	end := strings.Repeat("▬", int((1-f)*length))

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(start) + " 🚗 " + end
}

func (m *Model) RenderTyper() string {
	var cars []string
	const WIDTH int = 81

	for i, client := range m.clientsInLobby {
		// silly fixer
		if client.prog == nil {
			f := 0
			client.prog = &f
		}

		car := lipgloss.NewStyle().PaddingLeft(23).Render(renderCar(colors[i%len(colors)], *client.prog))
		name := lipgloss.NewStyle().Width(22).PaddingLeft(1).Render(client.getName())

		cars = append(cars, lipgloss.JoinHorizontal(lipgloss.Center, car, name))
	}

	carGloss := lipgloss.JoinVertical(lipgloss.Center, cars...)

	s := lipgloss.JoinVertical(lipgloss.Center, carGloss, lipgloss.
		NewStyle().
		Width(70).
		Foreground(lipgloss.Color("253")).
		Padding(4, 6, 0, 5).
		Render(m.renderText()),
	)

	return s

}

func (m *Model) renderText() string {

	v := m.typingInfo

	start := v.text[0:v.correctCharacters]
	typo := v.text[v.correctCharacters : v.correctCharacters+v.typoCharacters]
	end := v.text[v.correctCharacters+v.typoCharacters:]

	start = lipgloss.
		NewStyle().
		// Foreground(lipgloss.Color("253")).
		// Bold(true).
		Render(start)
	typo = lipgloss.
		NewStyle().
		Background(lipgloss.Color("9")).
		Foreground(lipgloss.Color("253")).
		Render(typo)
	end = lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(end)

	return lipgloss.JoinHorizontal(lipgloss.Center, start, typo, end)
}
