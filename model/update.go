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
	// listen for spinner tick
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	// listen for updated jq results
	case jq.ParseQueryResult:
		m.isLoading = false
		if msg.Err != nil {
			m.lastError = msg.Err
		} else {
			m.lastError = nil
			m.lastSuccessfulResult = msg.Result
		}
	}

	// handle text input changes
	m.textinput, cmd = m.textinput.Update(msg)
	query := util.SanitizeQuery(m.textinput.Value(), m.textinput.Placeholder)

	// if query changed, trigger new parsing of jq
	if query != m.lastQuery {
		m.lastQuery = query
		m.isLoading = true

		return m, func() tea.Msg {
			return jq.ParseQuery(m.data, query)
		}
	}

	return m, cmd
}
