package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/samsullivan/jqless/model"
)

func main() {
	// get file, expected to be JSON data
	file, err := getFile()
	if err != nil {
		panic(err)
	}

	// create new model with JSON file
	m, err := model.New(file)
	if err != nil {
		panic(err)
	}

	// start bubbletea program
	if _, err := tea.NewProgram(m).Run(); err != nil {
		panic(err)
	}
}

// getFile opens the first command line argument or piped data.
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

		if stat.Mode()&os.ModeNamedPipe != 0 && stat.Size() != 0 {
			file = os.Stdin
		}
	}

	return file, nil
}
