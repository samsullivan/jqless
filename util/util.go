package util

import "strings"

// SanitizeQuery removes leading and trailing spaces from input string.
// If string is empty, uses default value.
func SanitizeQuery(input string, defaultValue string) string {
	output := strings.TrimSpace(input)
	if output == "" {
		output = defaultValue
	}
	return output
}
