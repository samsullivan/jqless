package model

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	output := make([]string, 0, 4)

	output = append(output, strings.Join(m.lastResults, "\n"))
	if m.lastError != nil {
		output = append(output, fmt.Sprintf("error: %s", m.lastError))
	}

	output = append(output, m.spinner.View())
	output = append(output, m.textinput.View())
	output = append(output, "press ctrl+c to quit")

	return strings.Join(output, "\n\n")
}
