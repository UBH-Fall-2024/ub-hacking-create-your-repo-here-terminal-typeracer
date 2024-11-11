package tui

import "github.com/charmbracelet/lipgloss"

type Style struct {
	buttonStyle     lipgloss.Style
	errorStyle      lipgloss.Style
	connectingStyle lipgloss.Style
	notiStyle       lipgloss.Style
}

func NewStyle(renderer *lipgloss.Renderer) Style {
	purple := lipgloss.Color("57")
	errorGrey := lipgloss.Color("237")
	connectingPink := lipgloss.Color("219")
	notiPurple := lipgloss.Color("98")

	buttonStyle := lipgloss.NewStyle().
		Background(purple).
		Padding(0, 3).
		Underline(true)

	border := lipgloss.NewStyle().Padding(1, 2)

	return Style{
		buttonStyle:     buttonStyle,
		errorStyle:      border.Background(errorGrey),
		connectingStyle: border.Background(connectingPink),
		notiStyle: border.
			Border(lipgloss.RoundedBorder()).
			Foreground(notiPurple).
			MarginBottom(3),
	}
}
