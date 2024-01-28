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

// inputKeys contains the key binding & help text when input is focused.
var inputKeys = keyBindings{
	ViewportNavigation: getViewportNavigationKeyBinding(nil),
	Extract:            getExtractKeyBinding(true),
	SwitchFocus:        getSwitchFocusKeyBinding(nil),
	Quit:               getQuitKeyBinding(),
}

// inputKeys contains the key binding & help text when viewport is focused.
var viewportKeys = keyBindings{
	ViewportNavigation: getViewportNavigationKeyBinding([][]string{
		{"j", "k"},
		{"f", "b"},
		{"d", "u"},
	}),
	Extract:     getExtractKeyBinding(false),
	SwitchFocus: getSwitchFocusKeyBinding(util.Ptr("edit query")),
	Quit:        getQuitKeyBinding(),
}

// getViewportNavigationKeyBinding shows extra vim scrollable shortcuts.
func getViewportNavigationKeyBinding(extraHelpKeySets [][]string) key.Binding {
	keys := make([]string, 0, (len(extraHelpKeySets)*2)+1)
	keys = append(keys, "down", "up")

	helpKeys := make([]string, 0, len(extraHelpKeySets)+1)
	helpKeys = append(helpKeys, "↓/↑")

	for _, extraHelpKeySet := range extraHelpKeySets {
		keys = append(keys, extraHelpKeySet...)
		helpKeys = append(helpKeys, strings.Join(extraHelpKeySet, "/"))
	}

	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(strings.Join(helpKeys, "·"), "scroll output"),
	)
}

// getExtractKeyBinding allows optional enforcing of ctrl keypress.
func getExtractKeyBinding(requiresCtrl bool) key.Binding {
	keys := []string{"ctrl+x"}
	if !requiresCtrl {
		keys = append(keys, "x")
	}
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[len(keys)-1], "extract (to clipboard)"),
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

// getQuitKeyBinding has no options.
func getQuitKeyBinding() key.Binding {
	return key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	)
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k keyBindings) ShortHelp() []key.Binding {
	return []key.Binding{k.ViewportNavigation, k.Extract, k.SwitchFocus, k.Quit}
}

// FullHelp returns keybindings for the expanded help view; not implemented.
func (k keyBindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

// handleKeyMsg is used by Update() when a KeyMsg is received.
// Any non-nil command returned should be expected to be immediately
// passed to bubbletea, without being processed further by Update().
func (m model) handleKeyMsg(msg tea.KeyMsg) (model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, inputKeys.ViewportNavigation, viewportKeys.ViewportNavigation):
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	case key.Matches(msg, inputKeys.Extract, viewportKeys.Extract):
		cmd = func() tea.Msg {
			if err := util.WriteClipboard([]byte(m.viewportContents())); err != nil {
				return message.NewFatalError(err)
			}

			// TODO: success indication (flash help text green?)
			return nil
		}
		return m, cmd
	case key.Matches(msg, inputKeys.SwitchFocus, viewportKeys.SwitchFocus):
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
	case key.Matches(msg, inputKeys.Quit, viewportKeys.Quit):
		cmd = tea.Quit
		return m, cmd
	}

	return m, cmd
}
