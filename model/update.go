package model

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/samsullivan/jqless/jq"
	"github.com/samsullivan/jqless/message"
	"github.com/samsullivan/jqless/util"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		// some messages need to return a single command immediately
		cmd tea.Cmd
		// if not, append to list of cmds to be returned as a batch at the end
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	// listen for errors
	case message.FatalError:
		// TODO: better error output
		fmt.Printf("FatalError: %s\n\n", msg.Error())

		cmd = tea.Quit
		return m, cmd
	// listen to keypresses
	case tea.KeyMsg:
		m, cmd = m.handleKeyMsg(msg)
		if cmd != nil {
			return m, cmd
		}
	// listen for spinner tick
	case spinner.TickMsg:
		if !m.isLoading {
			// stop spinner if no longer loading
			return m, cmd
		}
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	// listen for window resizing
	case tea.WindowSizeMsg:
		// see https://github.com/charmbracelet/bubbletea/blob/master/examples/pager/main.go
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.viewportReady {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.SetContent(m.viewportContents())

			m.viewportReady = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
	// listen for parsed JSON file
	case message.ParsedFile:
		m.data = msg.Data()
	// listen for updated jq results
	case message.QueryResult:
		m.isLoading = false
		if msg.Failed() {
			m.lastError = msg.Error()
		} else {
			m.lastError = nil
			m.lastResults = msg.Results()
		}
	}

	// handle viewport changes
	m.viewport.SetContent(m.viewportContents())
	if m.currentFocus == focusViewport {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	// handle text input changes
	if m.currentFocus == focusInput {
		m.textinput, cmd = m.textinput.Update(msg)
		cmds = append(cmds, cmd)
	}

	// skip jq-related processing if file not processed into data yet
	if m.data == nil {
		// TODO: timeout
	} else {
		// if query changed, trigger new parsing of jq
		query := util.SanitizeQuery(m.textinput.Value(), m.textinput.Placeholder)
		if query != m.lastQuery {
			m.lastQuery = query
			m.isLoading = true

			// restart spinner in addition to triggering jq
			cmds = append(cmds, m.spinner.Tick)
			cmds = append(cmds, func() tea.Msg {
				return jq.Query(m.data, query, &jq.Options{
					Compact: m.compactOutput,
					Raw:     m.rawOutput,
				})
			})
		}
	}

	return m, tea.Batch(cmds...)
}
