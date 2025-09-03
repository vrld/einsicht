package ui

import (

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/k3a/html2text"
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
	CycleBodyType    key.Binding
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
	bodyToDisplay  string
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
			CycleBodyType:    key.NewBinding(key.WithKeys("tab", "b")),
		},
		Email:         email,
		OpenCommand:   viper.GetString("command"),
		inputMode:     uiModeReadBody,
		bodyToDisplay: viper.GetString("body"),
	}

	model.viewportHeader = viewport.New(1, 1)
	model.viewportBody = viewport.New(1, 1)
	model.setDimensions(1, 1)

	return model
}

func (m *Model) setDimensions(width, height int) {
	m.width = width
	m.height = height

	m.viewportHeader.Width = width - 2
	m.viewportHeader.SetContent(
		renderHeadersForDisplay(m.Email.HeaderDisplayAll(), m.width-4, true),
	)
	m.viewportHeader.Height = min(viewportHeaderHeight, m.viewportHeader.TotalLineCount())

	m.updateBodyViewport()
}

func (m *Model) setInputState(state int) {
	m.inputMode = state
	m.viewportBody.Height = m.computeViewportBodyHeight()
}

func (m *Model) updateBodyViewport() {
	m.viewportBody.Width = m.width - 2
	m.viewportBody.Height = m.computeViewportBodyHeight()

	content := ""
	switch m.bodyToDisplay {
	case "plain":
		content = m.Email.Text
	case "html":
		content = html2text.HTML2Text(m.Email.HTML)
	}

	limitWidth := lipgloss.NewStyle().Width(m.viewportBody.Width - 2)
	m.viewportBody.SetContent(limitWidth.Render(content))
}
