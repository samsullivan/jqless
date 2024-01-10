package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/itchyny/gojq"
)

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		panic("jqless expects piped data")
	}

	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var m model
	err = json.Unmarshal(b, &m.data)
	if err != nil {
		panic(err)
	}

	ti := textinput.New()
	ti.Placeholder = "."
	ti.Focus()
	m.input = ti

	query, err := gojq.Parse(ti.Placeholder)
	if err != nil {
		panic(err)
	}
	m.query = query

	if _, err := tea.NewProgram(m).Run(); err != nil {
		panic(err)
	}
}

type model struct {
	data  interface{}
	input textinput.Model
	query *gojq.Query
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.Type {
		case tea.KeyCtrlC, tea.KeyEscape, tea.KeyEnter:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	query, _ := gojq.Parse(m.input.Value()) // TODO: error handling
	if query != nil {
		m.query = query
	}

	return m, cmd
}

func (m model) View() string {
	var pieces []string

	iter := m.query.Run(m.data)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			panic(err)
		}

		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			panic(err)
		}

		pieces = append(pieces, string(b))
	}

	return fmt.Sprintf("%v\n\n%s\n\npress ctrl+c to quit", strings.Join(pieces, "\n"), m.input.View())
}
