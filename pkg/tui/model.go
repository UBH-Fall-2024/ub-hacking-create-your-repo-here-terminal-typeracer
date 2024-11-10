package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

type Model struct {
	Renderer *lipgloss.Renderer
	Width    int
	Height   int
}

func NewModel(renderer *lipgloss.Renderer, pty *ssh.Pty) *Model {
	return &Model{
		Renderer: renderer,
		Width:    pty.Window.Width,
		Height:   pty.Window.Height,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m *Model) View() string {
	return "HII!!!!!"
}
