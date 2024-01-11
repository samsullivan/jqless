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
	var b []byte
	if len(os.Args) > 1 {
		var err error
		b, err = os.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
	} else {
		stat, err := os.Stdin.Stat()
		if err != nil {
			panic(err)
		}

		if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
			panic("jqless expects piped data if filename not included as an argument")
		}

		b, err = io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
	}

	var m model
	err := json.Unmarshal(b, &m.data)
	if err != nil {
		panic(err)
	}

	ti := textinput.New()
	ti.Placeholder = "."
	ti.Focus()
	m.input = ti

	if _, err := tea.NewProgram(m).Run(); err != nil {
		panic(err)
	}
}

type model struct {
	data  interface{}
	input textinput.Model

	lastError            error
	lastQuery            string
	lastSuccessfulResult []string
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

	input := strings.TrimSpace(m.input.Value())
	if input == "" {
		input = m.input.Placeholder
	}

	if input != m.lastQuery {
		m.lastQuery = input

		// reset last error before trying again
		// intentionally leaving last successful state, in case this attempt fails
		m.lastError = nil

		query, err := gojq.Parse(input)
		if err != nil {
			m.lastError = err
		} else {
			var result []string

			iter := query.Run(m.data)
			for {
				v, ok := iter.Next()
				if !ok {
					break
				}
				if err, ok := v.(error); ok {
					// TODO: handle more than one error
					m.lastError = err
					break
				}

				b, err := json.MarshalIndent(v, "", "  ")
				if err != nil {
					m.lastError = err
					break
				}

				result = append(result, string(b))
			}

			if m.lastError == nil {
				m.lastSuccessfulResult = result
			}
		}
	}

	return m, cmd
}

func (m model) View() string {
	output := make([]string, 0, 4)

	output = append(output, strings.Join(m.lastSuccessfulResult, "\n"))
	if m.lastError != nil {
		output = append(output, fmt.Sprintf("error: %s", m.lastError))
	}

	output = append(output, m.input.View())
	output = append(output, "press ctrl+c to quit")

	return strings.Join(output, "\n\n")
}
