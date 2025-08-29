package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"github.com/vrld/einsicht/internal"
)

func Run(email *internal.Email) error {
	zone.NewGlobal()
	defer zone.Close()

	model := InitialModel(email)

	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	_, err := p.Run()
	return err
}
