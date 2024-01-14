package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/samsullivan/jqless/util"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

func (m model) View() string {
	if !m.viewportReady || m.data == nil {
		return ""
	}

	return strings.Join([]string{
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
	}, "\n")
}

func (m model) headerView() string {
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
	info := infoStyle.Render(fmt.Sprintf("%s %3.f%%", m.spinner.View(), m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", util.Max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
