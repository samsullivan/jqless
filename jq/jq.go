package jq

import (
	"encoding/json"

	"github.com/itchyny/gojq"
	"github.com/samsullivan/jqless/message"
)

const DefaultQuery = "."

// Query takes unmarshalled JSON data and an input query.
// On success, returns a slice of strings from gojq result.
func Query(data interface{}, input string) message.QueryResult {
	query, err := gojq.Parse(input)
	if err != nil {
		return message.NewQueryError(err)
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
			return message.NewQueryError(err)
		}

		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return message.NewQueryError(err)
		}

		result = append(result, string(b))
	}

	return message.NewQueryResult(result)
}
