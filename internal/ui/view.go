package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"

	"github.com/vrld/einsicht/internal"
)

const (
	borderColorNormal   = lipgloss.Color("8")
	borderColorSelected = lipgloss.Color("4")
)

const viewportHeaderHeight = 10

var (
	styleHeaderKey   = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	styleHeaderValue = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)

	styleButton       = lipgloss.NewStyle().Background(lipgloss.Color("4")).Foreground(lipgloss.Color("0")).Padding(0, 1).Bold(true)
	styleCancelButton = styleButton.Background(lipgloss.Color("9"))
	styleInfo         = styleButton.Background(lipgloss.Color("6"))
)

func (m *Model) borderColor(state int) lipgloss.Color {
	if m.inputMode == state {
		return borderColorSelected
	}
	return borderColorNormal
}

func (m Model) View() string {
	cardHeaders := m.renderHeaders()

	var cardAttachments string
	if len(m.Email.Attachments) > 0 {
		cardAttachments = m.renderAttachments() + "\n"
	}

	bottom := m.renderBottom()

	cardBody := m.renderBody()

	// NOTE: cardAttachments already includes the \n if there are attachments
	return zone.Scan(fmt.Sprint(cardHeaders, "\n", cardBody, "\n", cardAttachments, bottom))
}

func (m *Model) renderHeightHeaders() int {
	if m.inputMode == uiModeInspectHeader {
		return viewportHeaderHeight + 2
	}
	return len(m.Email.HeaderDisplayCanonical()) + 2
}

func (m *Model) renderHeaders() string {
	color := m.borderColor(uiModeInspectHeader)

	headerList := ""
	if m.inputMode == uiModeInspectHeader {
		headerList = m.viewportHeader.View()
	} else {
		headerList = renderHeadersForDisplay(m.Email.HeaderDisplayCanonical(), m.width-2)
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, true, true).
		BorderForeground(color).
		Padding(0, 1).
		Width(m.width - 2)

	return fmt.Sprint(
		renderTopLineWithTitle("header", m.width, color),
		"\n",
		style.Render(headerList),
	)
}

func renderHeadersForDisplay(headers []internal.HeaderDisplay, width int) string {
	const padding = 2

	headerDisplayWidth := 0
	for _, h := range headers {
		headerDisplayWidth = max(headerDisplayWidth, len(h.Key))
	}

	valueDisplayWidth := width - headerDisplayWidth - padding

	lines := []string{}
	for _, header := range headers {
		lines = append(lines, fmt.Sprint(
			styleHeaderKey.Width(headerDisplayWidth + padding).Render(header.Key+":"),
			styleHeaderValue.Width(valueDisplayWidth).Render(header.Value),
		))
	}

	return strings.Join(lines, "\n")
}

func (m *Model) renderHeightAttachments() int {
	count := len(m.Email.Attachments)
	if count == 0 {
		return 0
	}

	padding := 2
	if m.inputMode == uiModeSelectAttachment {
		padding += len(m.Email.Attachments) - 1
	}

	return count + padding
}

func (m *Model) renderAttachments() string {
	lines := []string{}
	for i, a := range m.Email.Attachments {
		line := ""
		if m.inputMode == uiModeSelectAttachment {
			newline := ""
			if i > 0 {
				newline = "\n"
			}
			selectionRune := attachmentIndexToRune(i)
			prefix := fmt.Sprintf("⟨%c⟩", selectionRune)
			line = fmt.Sprint(newline, zone.Mark(
				fmt.Sprintf("attachment:%c", selectionRune),
				styleButton.Render(prefix, a.Filename, a.ContentType, internal.HumanReadableFileSize(len(a.Content))),
			))
		} else {
			line = fmt.Sprint("• ", a.Filename, " ", a.ContentType, " ", internal.HumanReadableFileSize(len(a.Content)))
		}
		lines = append(lines, line)
	}

	color := m.borderColor(uiModeSelectAttachment)

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, true, true).
		BorderForeground(color).
		Padding(0, 1).
		Width(m.width - 2)

	return fmt.Sprint(
		renderTopLineWithTitle("attachments", m.width, color),
		"\n",
		style.Render(strings.Join(lines, "\n")),
	)
}

func (m *Model) renderHeightBottom() int {
	return 1
}

func (m *Model) renderBottom() string {
	switch m.inputMode {
	case uiModeSelectAttachment:
		return m.renderButtons(
			[]string{styleInfo.Render("select attachment to save")},
			[]string{
				zone.Mark("cancelSelect", styleCancelButton.Render("cancel <esc>")),
			},
		)

	case uiModeInspectHeader:
		return m.renderButtons(
			[]string{styleInfo.Render("inspect headers")},
			[]string{
				zone.Mark("cancelSelect", styleCancelButton.Render("cancel <esc>")),
			},
		)
	}

	// main mode

	buttonsLeft := []string{
		zone.Mark("openHTML", styleButton.Render("[o]pen HTML")),
	}

	if len(m.Email.Attachments) > 0 {
		buttonsLeft = append(buttonsLeft, zone.Mark("attach", styleButton.Render("save [a]ttachment")))
	}

	buttonsLeft = append(buttonsLeft, []string{
		zone.Mark("headers", styleButton.Render("inspect [h]eaders")),
		zone.Mark("save", styleButton.Render("[s]ave body")),
	}...)

	return m.renderButtons(
		buttonsLeft,
		[]string{
			zone.Mark("quit", styleCancelButton.Render("[q]uit")),
		},
	)
}

func (m *Model) renderButtons(left, right []string) string {
	leftJoined := strings.Join(left, " ")
	rightJoined := strings.Join(right, " ")

	spacer := ""
	spacerWidth := m.width - (lipgloss.Width(leftJoined) + lipgloss.Width(rightJoined)) - 2
	if spacerWidth > 0 {
		spacer = strings.Repeat(" ", spacerWidth)
	}

	return fmt.Sprint(" ", leftJoined, spacer, rightJoined)
}

func (m *Model) renderBody() string {
	color := m.borderColor(uiModeReadBody)

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, true, true).
		BorderForeground(lipgloss.Color(color)).
		Padding(0, 1).
		Width(m.width - 2)

	return fmt.Sprint(
		renderTopLineWithTitle("body", m.width, lipgloss.Color(color)),
		"\n",
		style.Height(m.viewportBody.Height-1).Render(m.viewportBody.View()),
	)
}

func (m *Model) computeViewportBodyHeight() int {
	return m.height - m.renderHeightAttachments() - m.renderHeightHeaders() - m.renderHeightBottom() - 2
}

func renderTopLineWithTitle(title string, width int, color lipgloss.Color) string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	switch width {
	case 0:
		return ""
	case 1:
		return borderStyle.Render("╌")
	case 2:
		return borderStyle.Render("├┤")
	}

	titleWidth := lipgloss.Width(title) + 6
	if width < titleWidth {
		return borderStyle.Render(fmt.Sprint("┌", strings.Repeat("─", width-2), "┐"))
	}

	return borderStyle.Render(
		fmt.Sprint(
			"┌",
			strings.Repeat("─", width-titleWidth),
			"─🮤",
			title,
			"🮥─┐",
		),
	)
}

func attachmentIndexToRune(index int) rune {
	if index >= 0 && index <= 9 {
		return rune('0' + ((index + 1) % 10))
	}

	if index >= 10 && index <= 35 {
		return rune('a' + index - 10)
	}

	return 'ẞ'
}
