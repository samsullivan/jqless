package main

import (
	"context"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v3"

	"github.com/samsullivan/jqless/model"
)

var version = "0.0.0"

// main registers cli flags and, if unused, starts the bubbletea program.
func main() {
	cmd := &cli.Command{
		Name:    "jqless",
		Usage:   "combining jq and less for real-time JSON parsing",
		Version: version,
		UsageText: strings.Join([]string{
			"jqless [path/to/file.json]",
			"cat [path/to/file.json] | jqless",
		}, "\n"),
		HideHelpCommand: true,
		Action: func(context.Context, *cli.Command) error {
			file, err := getFile()
			if err != nil {
				return err
			}

			m, err := model.New(file)
			if err != nil {
				return err
			}

			if _, err := tea.NewProgram(m).Run(); err != nil {
				return err
			}

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
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
