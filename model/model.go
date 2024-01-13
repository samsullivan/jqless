package model

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/itchyny/gojq"
)

type model struct {
	data  gojq.PreparedData
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

	// prepare input JSON for concurrent gojq execution while querying
	m.data = gojq.PrepareData(m.data)

	// configure text input
	ti := textinput.New()
	ti.Placeholder = "."
	ti.Focus()
	m.input = ti

	return &m, nil
}
