package tui

import "github.com/charmbracelet/lipgloss"

type Style struct {
	buttonStyle     lipgloss.Style
	errorStyle      lipgloss.Style
	connectingStyle lipgloss.Style
}

func NewStyle(renderer *lipgloss.Renderer) Style {
	purple := lipgloss.Color("57")
	errorGrey := lipgloss.Color("237")
	connectingOrange := lipgloss.Color("202")

	buttonStyle := lipgloss.NewStyle().
		Background(purple).
		Padding(0, 4)

	border := lipgloss.NewStyle().Padding(1, 2)

	return Style{
		buttonStyle:     buttonStyle,
		errorStyle:      border.Background(errorGrey),
		connectingStyle: border.Background(connectingOrange),
	}
}
