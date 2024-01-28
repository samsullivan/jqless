package model

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/samsullivan/jqless/util"
)

var (
	leftBoxStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}

	rightBoxStyleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return leftBoxStyle().Copy().BorderStyle(b)
	}
)

func (m model) View() string {
	if m.file == nil {
		// FatalError will be returned by parseFile() command triggered by Init().
		return ""
	}

	if !m.viewportReady || m.data == nil {
		return m.footerView()
	}

	return strings.Join([]string{
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
	}, "\n")
}

func (m model) headerView() string {
	titleStyle := leftBoxStyle()

	var borderColor lipgloss.TerminalColor = lipgloss.NoColor{}
	if m.lastError != nil {
		borderColor = lipgloss.Color("#CF2222")
	}
	titleStyle.BorderForeground(borderColor)

	title := titleStyle.Render(m.textinput.View())
	line := strings.Repeat("─", util.Max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) viewportContents() string {
	return strings.Join(m.lastResults, "\n")
}

func (m model) footerView() string {
	k := keys
	if m.currentFocus == focusViewport {
		k.ViewportNavigation = getViewportNavigationKeyBinding([]string{
			"j/k",
			"f/b",
			"d/u",
		})
	}
	help := leftBoxStyle().Render(m.help.View(k))

	infoItems := make([]string, 1, 2)
	infoItems[0] = m.spinner.View()
	if m.viewport.TotalLineCount() > m.viewport.VisibleLineCount() {
		infoItems = append(infoItems, fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	}

	slices.Reverse(infoItems)
	info := rightBoxStyleStyle().Render(strings.Join(infoItems, " "))

	contentWidth := lipgloss.Width(help) + lipgloss.Width(info)
	line := strings.Repeat("─", util.Max(0, m.viewport.Width-contentWidth))
	return lipgloss.JoinHorizontal(lipgloss.Center, help, line, info)
}
