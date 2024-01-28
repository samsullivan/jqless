package model

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/samsullivan/jqless/message"
	"github.com/samsullivan/jqless/util"
)

// keyBindings defines the available key bindings.
type keyBindings struct {
	ViewportNavigation key.Binding
	Extract            key.Binding
	SwitchFocus        key.Binding
	Quit               key.Binding
}

// Validate that keyBindings satisfies help.KeyMap interface.
var _ help.KeyMap = (*keyBindings)(nil)

// getViewportNavigationKeyBinding shows extra vim scrollable shortcuts.
func getViewportNavigationKeyBinding(extraHelpKeys []string) key.Binding {
	helpKeys := make([]string, 0, len(extraHelpKeys)+1)
	helpKeys = append(helpKeys, "↓/↑")
	helpKeys = append(helpKeys, extraHelpKeys...)
	return key.NewBinding(
		key.WithKeys("down", "up"), // always down/up, regardless of extraHelpKeys
		key.WithHelp(strings.Join(helpKeys, "·"), "scroll output"),
	)
}

// getSwitchFocusKeyBinding allows overriding the help text.
func getSwitchFocusKeyBinding(customHelpText *string) key.Binding {
	helpText := "more options"
	if customHelpText != nil {
		helpText = *customHelpText
	}
	return key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("⇥", helpText),
	)
}

// keys contains the actual key bindings as well as the related help text.
var keys = keyBindings{
	ViewportNavigation: getViewportNavigationKeyBinding(nil),
	Extract: key.NewBinding(
		key.WithKeys("ctrl+x"),
		key.WithHelp("ctrl+x", "extract (to clipboard)"),
	),
	SwitchFocus: getSwitchFocusKeyBinding(nil),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k keyBindings) ShortHelp() []key.Binding {
	return []key.Binding{k.ViewportNavigation, k.Extract, k.SwitchFocus, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k keyBindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ViewportNavigation, k.Extract, k.SwitchFocus, k.Quit},
		// TODO: additional actions
	}
}

// handleKeyMsg is used by Update() when a KeyMsg is received.
// Any non-nil command returned should be expected to be immediately
// passed to bubbletea, without being processed further by Update().
func (m model) handleKeyMsg(msg tea.KeyMsg) (model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, keys.ViewportNavigation):
		// this is only needed for ↓/↑ navigation; when focus on viewport,
		// the main Update() will handle all viewport update messages
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	case key.Matches(msg, keys.Extract):
		cmd = func() tea.Msg {
			if err := util.WriteClipboard([]byte(m.viewportContents())); err != nil {
				return message.NewFatalError(err)
			}

			// TODO: success indication (flash help text green?)
			return nil
		}
		return m, cmd
	case key.Matches(msg, keys.SwitchFocus):
		switch m.currentFocus {
		case focusInput:
			m.currentFocus = focusViewport
			m.textinput.Cursor.Blink = true
		case focusViewport:
			m.currentFocus = focusInput
			m.textinput.Cursor.Blink = false
			cmd = textinput.Blink
		}
		return m, cmd
	case key.Matches(msg, keys.Quit):
		cmd = tea.Quit
		return m, cmd
	}

	return m, cmd
}
