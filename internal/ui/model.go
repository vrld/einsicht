package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/spf13/viper"
	"github.com/vrld/einsicht/internal"
)

type KeyMap struct {
	OpenHTML         key.Binding
	SelectAttachment key.Binding
	InspectHeaders   key.Binding
	Save             key.Binding
	Quit             key.Binding
	SetDefaultMode   key.Binding
}

const (
	uiModeReadBody = iota
	uiModeSelectAttachment
	uiModeInspectHeader
)

type Model struct {
	width, height  int
	KeyMap         KeyMap
	Email          *internal.Email
	OpenCommand    string
	viewportBody   viewport.Model
	viewportHeader viewport.Model
	inputMode      int
}

func InitialModel(email *internal.Email) Model {
	model := Model{
		KeyMap: KeyMap{
			OpenHTML:         key.NewBinding(key.WithKeys("o")),
			SelectAttachment: key.NewBinding(key.WithKeys("a")),
			InspectHeaders:   key.NewBinding(key.WithKeys("h")),
			Save:             key.NewBinding(key.WithKeys("s")),
			Quit:             key.NewBinding(key.WithKeys("q", "ctrl+c")),
			SetDefaultMode:   key.NewBinding(key.WithKeys("escape", "q")),
		},
		Email:     email,
		OpenCommand: viper.GetString("command"),
		inputMode: uiModeReadBody,
	}

	model.viewportHeader = viewport.New(1, 1)
	model.viewportHeader.Height = viewportHeaderHeight
	model.viewportBody = viewport.New(1, 1)
	model.setDimensions(1, 1)
	model.UpdateEmailDisplay()

	return model
}

func (m *Model) setDimensions(width, height int) {
	m.width = width
	m.height = height

	m.viewportHeader.Width = width - 2
	m.viewportHeader.SetContent(
		renderHeadersForDisplay(m.Email.HeaderDisplayAll(), m.viewportHeader.Width),
	)
	m.viewportBody.Width = width - 2
	m.viewportBody.Height = m.computeViewportBodyHeight()
}

func (m *Model) setInputState(state int) {
	m.inputMode = state
	m.viewportBody.Height = m.computeViewportBodyHeight()
}

func (m *Model) UpdateEmailDisplay() {
	m.viewportBody.SetContent(string(m.Email.Text))
	m.viewportBody.Height = m.computeViewportBodyHeight()
}
