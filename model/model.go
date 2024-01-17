package model

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/samsullivan/jqless/jq"
	"github.com/samsullivan/jqless/message"
)

type model struct {
	viewportReady bool

	// related to JSON user input
	file *os.File
	data interface{}

	// bubbletea components
	viewport  viewport.Model
	help      help.Model
	textinput textinput.Model
	spinner   spinner.Model

	// state of jq querying
	isLoading   bool
	lastError   error
	lastQuery   string
	lastResults []string
}

// New takes an open file and returns a model for use by bubbletea.
// In order to show the spinner immediately, for larger JSON payloads,
// The file stream isn't consumed or unmarshalled into JSON yet.
func New(file *os.File) (*model, error) {
	m := model{file: file}

	// configure help
	m.help = help.New()

	// configure text input
	m.textinput = textinput.New()
	m.textinput.Placeholder = jq.DefaultQuery
	m.textinput.Focus()

	// configure loading spinner
	m.spinner = spinner.New()

	return &m, nil
}

// parseFile returns a command for reading the input file into unmarshalled JSON data
func (m *model) parseFile() tea.Cmd {
	return func() tea.Msg {
		var data interface{}

		// verify file exists
		if m.file == nil {
			return message.NewFatalError(
				errors.New("no data passed to jqless"),
			)
		}

		// close file when done reading
		defer m.file.Close()

		// read entire file
		b, err := io.ReadAll(m.file)
		if err != nil {
			return message.NewFatalError(err)
		}

		// unmarshal to data interface
		err = json.Unmarshal(b, &data)
		if err != nil {
			return message.NewFatalError(err)
		}

		return message.NewParsedFile(data)
	}
}
