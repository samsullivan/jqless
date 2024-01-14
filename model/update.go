package model

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/samsullivan/jqless/jq"
	"github.com/samsullivan/jqless/util"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// listen to keypresses
	case tea.KeyMsg:
		switch msg.Type {
		// exit keys
		case tea.KeyCtrlC, tea.KeyEscape, tea.KeyEnter:
			cmd = tea.Quit
			return m, cmd
		}
	// listen for loading spinner
	case spinner.TickMsg:
		m.loading, cmd = m.loading.Update(msg)
		return m, cmd
	// listen for updated jq results
	case jq.ParseQueryResult:
		if msg.Err != nil {
			m.lastError = msg.Err
		} else {
			m.lastError = nil
			m.lastSuccessfulResult = msg.Result
		}
	}

	// handle text input changes
	m.input, cmd = m.input.Update(msg)
	query := util.SanitizeQuery(m.input.Value(), m.input.Placeholder)

	// if query changed, trigger new parsing of jq
	if query != m.lastQuery {
		m.lastQuery = query

		return m, func() tea.Msg {
			return jq.ParseQuery(m.data, query)
		}
	}

	return m, cmd
}
