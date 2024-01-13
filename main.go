package main

import (
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/samsullivan/jqless/model"
)

func main() {
	var b []byte

	if len(os.Args) > 1 {
		// if command line arguments included, attempt to open local file
		var err error
		b, err = os.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}
	} else {
		// otherwise, piped data is expected
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

	// create new model with input JSON
	m, err := model.New(b)
	if err != nil {
		panic(err)
	}

	// start bubbletea program
	if _, err := tea.NewProgram(m).Run(); err != nil {
		panic(err)
	}
}
