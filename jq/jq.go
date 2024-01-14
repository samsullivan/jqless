package jq

import (
	"encoding/json"

	"github.com/itchyny/gojq"
)

const DefaultQuery = "."

type ParseQueryResult struct {
	Result []string
	Err    error
}

// ParseQuery takes unmarshalled JSON data and an input query.
// On success, returns a slice of strings from gojq result.
func ParseQuery(data interface{}, input string) ParseQueryResult {
	query, err := gojq.Parse(input)
	if err != nil {
		return ParseQueryResult{
			Err: err,
		}
	}

	var result []string

	iter := query.Run(data)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			// TODO: handle more than one error
			return ParseQueryResult{
				Err: err,
			}
		}

		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return ParseQueryResult{
				Err: err,
			}
		}

		result = append(result, string(b))
	}

	return ParseQueryResult{
		Result: result,
	}
}
