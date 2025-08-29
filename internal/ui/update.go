package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		panic(msg)

	case OpenHTMLMsg:
		// TODO

	case SaveAttachmentMsg:
		// TODO

	case OpenAttachmentMsg:
		// TODO

	case SaveBodyMsg:
		// TODO

	case tea.WindowSizeMsg:
		m.setDimensions(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch m.inputMode {
		case uiModeReadBody:
			return m, m.handleReadBodyKeys(msg)

		case uiModeSelectAttachment:
			return m, m.handleSelectAttachmentKeys(msg)

		case uiModeInspectHeader:
			return m, m.handleInspectHeaderKeys(msg)
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
			switch m.inputMode {
			case uiModeReadBody:
				if zone.Get("openHtml").InBounds(msg) {
					return m, send(OpenHTMLMsg{})
				} else if zone.Get("attach").InBounds(msg) && len(m.Email.Attachments) > 0 {
					m.setInputState(uiModeSelectAttachment)
				} else if zone.Get("headers").InBounds(msg) {
					m.setInputState(uiModeInspectHeader)
				} else if zone.Get("save").InBounds(msg) {
					return m, send(SaveBodyMsg{})
				} else if zone.Get("quit").InBounds(msg) {
					return m, tea.Quit
				}

			case uiModeSelectAttachment:
				if zone.Get("cancelSelect").InBounds(msg) {
					m.setInputState(uiModeReadBody)
				}
				for i := range m.Email.Attachments {
					zoneName := fmt.Sprintf("attachment:%c", attachmentIndexToRune(i))
					if zone.Get(zoneName).InBounds(msg) {
						return m, send(SaveAttachmentMsg{i})
					}
				}
			}
			return m, nil
		}

	}

	var cmd tea.Cmd
	m.viewportBody, cmd = m.viewportBody.Update(msg)

	return m, cmd
}

func (m *Model) handleReadBodyKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.KeyMap.OpenHTML):
		return send(OpenHTMLMsg{})

	case key.Matches(msg, m.KeyMap.SelectAttachment) && len(m.Email.Attachments) > 0:
		m.setInputState(uiModeSelectAttachment)

	case key.Matches(msg, m.KeyMap.InspectHeaders):
		m.setInputState(uiModeInspectHeader)

	case key.Matches(msg, m.KeyMap.Save):
		return send(SaveBodyMsg{})

	case key.Matches(msg, m.KeyMap.Quit):
		return tea.Quit

	default:
		var cmdBody, cmdHeader tea.Cmd
		m.viewportBody, cmdBody = m.viewportBody.Update(msg)
		m.viewportHeader, cmdHeader = m.viewportHeader.Update(msg)
		return tea.Batch(cmdBody, cmdHeader)
	}

	return nil
}

func (m *Model) handleInspectHeaderKeys(msg tea.KeyMsg) (cmd tea.Cmd) {
	if key.Matches(msg, m.KeyMap.SetDefaultMode, m.KeyMap.InspectHeaders)  {
		m.setInputState(uiModeReadBody)
		return nil
	}

	m.viewportHeader, cmd = m.viewportHeader.Update(msg)
	return cmd
}

func (m *Model) handleSelectAttachmentKeys(msg tea.KeyMsg) tea.Cmd {
	if key.Matches(msg, m.KeyMap.SetDefaultMode) {
		m.setInputState(uiModeReadBody)
		return nil
	}

	attachmentIndex := runeToAttachmentIndex(msg)
	if attachmentIndex >= 0 && attachmentIndex < len(m.Email.Attachments) {
		return send(SaveAttachmentMsg{attachmentIndex})
	}

	return nil
}


func runeToAttachmentIndex(msg tea.KeyMsg) int {
	if len(msg.Runes) == 0 {
		return -1
	}

	rune := msg.Runes[0]
	// 1 -> 0, 2 -> 1, ..., 0 -> 9
	if rune >= '0' && rune <= '9' {
		return (int(rune-'0') + 9) % 10
	}

	// a -> 10, b -> 11, ...
	if rune >= 'a' && rune <= 'z' {
		return int(rune-'a') + 10
	}

	return -1
}

type OpenHTMLMsg struct{}
type SaveAttachmentMsg struct{Index int}
type OpenAttachmentMsg struct{Index int}
type SaveBodyMsg struct{}

func send(value tea.Msg) tea.Cmd {
	return func() tea.Msg { return value }
}
