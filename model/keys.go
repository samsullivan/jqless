package model

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/samsullivan/jqless/message"
	"github.com/samsullivan/jqless/util"
)

// keyBindings defines the available key bindings.
type keyBindings struct {
	Extract key.Binding
	Quit    key.Binding
}

// Validate that keyBindings satisfies help.KeyMap interface.
var _ help.KeyMap = (*keyBindings)(nil)

// keys contains the actual key bindings as well as the related help text.
var keys = keyBindings{
	Extract: key.NewBinding(
		key.WithKeys("ctrl+x"),
		key.WithHelp("ctrl+x", "extract (to clipboard)"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k keyBindings) ShortHelp() []key.Binding {
	return []key.Binding{k.Extract, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k keyBindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Extract, k.Quit},
		// TODO: additional actions
	}
}

// handleKeyMsg is used by Update() when a KeyMsg is received.
// Any non-nil command returned should be expected to be immediately
// passed to bubbletea, without being processed further by Update().
func (m model) handleKeyMsg(msg tea.KeyMsg) (model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, keys.Extract):
		cmd = func() tea.Msg {
			if err := util.WriteClipboard([]byte(m.viewportContents())); err != nil {
				return message.NewFatalError(err)
			}

			// TODO: success indication (flash help text green?)
			return nil
		}
		return m, cmd
	case key.Matches(msg, keys.Quit):
		cmd = tea.Quit
		return m, cmd
	}

	return m, cmd
}
