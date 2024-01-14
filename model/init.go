package model

import tea "github.com/charmbracelet/bubbletea"

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}
