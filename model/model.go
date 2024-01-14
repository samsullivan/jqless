package model

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"

	"github.com/samsullivan/jqless/jq"
)

type model struct {
	// data is unmarshalled JSON user input
	data interface{}

	// bubbletea components
	textinput textinput.Model
	spinner   spinner.Model

	// state of jq querying
	isLoading            bool
	lastError            error
	lastQuery            string
	lastSuccessfulResult []string
}

// New takes input JSON as a byte slice and returns a model for use by bubbletea.
func New(input []byte) (*model, error) {
	var m model

	// parse input JSON, either from local file or piped data
	err := json.Unmarshal(input, &m.data)
	if err != nil {
		panic(err)
	}

	// configure text input
	ti := textinput.New()
	ti.Placeholder = jq.DefaultQuery
	ti.Focus()
	m.textinput = ti

	// configure loading spinner
	spin := spinner.New()
	m.spinner = spin

	return &m, nil
}
