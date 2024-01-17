package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/samsullivan/jqless/model"
)

// main starts the bubbletea program
func main() {
	file, err := getFile()
	if err != nil {
		panic(err)
	}

	m, err := model.New(file)
	if err != nil {
		panic(err)
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		panic(err)
	}
}

// getFile returns a file descriptor, expected to contain JSON data.
func getFile() (file *os.File, err error) {
	if len(os.Args) > 1 {
		// if command line arguments included, attempt to open local file
		file, err = os.Open(os.Args[1])
		if err != nil {
			return nil, err
		}
	} else {
		// otherwise, piped data is expected
		stat, err := os.Stdin.Stat()
		if err != nil {
			return nil, err
		}

		if stat.Mode()&os.ModeNamedPipe != 0 {
			file = os.Stdin
		}
	}

	return file, nil
}
