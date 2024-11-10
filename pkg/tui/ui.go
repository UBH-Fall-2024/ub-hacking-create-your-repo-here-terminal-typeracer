package tui

import (
	"log"
	"strings"

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
			MarginRight(1)).ItemStyle(lipgloss.NewStyle().Width(22))

	return l.String()
}

func renderCar(color string, prog int) string {
	const length = 32

	f := prog / 100

	start := strings.Repeat("▬", f*length)
	end := strings.Repeat("▬", (1-f)*length)

	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(start) + " 🚗 " + end
}

func (m *Model) RenderTyper() string {
	var s string
	const WIDTH int = 81

	for i, client := range m.clientsInLobby {
		// silly fixer
		if client.prog == nil {
			f := 0
			client.prog = &f
		}

		car := lipgloss.NewStyle().PaddingLeft(23).Render(renderCar(colors[i%len(colors)], *client.prog))
		name := lipgloss.NewStyle().Width(22).PaddingRight(1).Render(client.getName())

		s += car + name + "\n"
	}

	s += lipgloss.NewStyle().Width(70).Padding(4, 6, 0, 5).Render(m.renderText())

	return s

}

func (m *Model) renderText() string {

	v := m.typingInfo

	start := v.text[0:v.correctCharacters]
	typo := v.text[v.correctCharacters : v.correctCharacters+v.typoCharacters]
	end := v.text[v.correctCharacters+v.typoCharacters:]

	s := start
	s += lipgloss.NewStyle().Background(lipgloss.Color("9")).Render(typo)
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(end)

	return s
}
