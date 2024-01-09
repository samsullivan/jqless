package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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

	if _, err := tea.NewProgram(m).Run(); err != nil {
		panic(err)
	}
}

type model struct {
	data interface{}
}

func (m model) Init() tea.Cmd {
	// no I/O for now
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.Type {
		case tea.KeyCtrlC, tea.KeyEscape, tea.KeyEnter:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	b, err := json.MarshalIndent(m.data, "", "\t")
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s\npress ctrl+c to quit", string(b))
}
