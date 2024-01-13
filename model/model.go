package model

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/textinput"
)

type model struct {
	data  interface{}
	input textinput.Model

	lastError            error
	lastQuery            string
	lastSuccessfulResult []string
}

func New(input []byte) (*model, error) {
	var m model

	// parse input JSON, either from local file or piped data
	err := json.Unmarshal(input, &m.data)
	if err != nil {
		panic(err)
	}

	// configure text input
	ti := textinput.New()
	ti.Placeholder = "."
	ti.Focus()
	m.input = ti

	return &m, nil
}
