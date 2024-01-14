package message

type QueryResult struct {
	results []string
	err     error
}

func NewQueryResult(results []string) QueryResult {
	return QueryResult{results: results}
}

func NewQueryError(err error) QueryResult {
	return QueryResult{err: err}
}

func (msg QueryResult) Results() []string {
	return msg.results
}

func (msg QueryResult) Error() error {
	return msg.err
}

func (msg QueryResult) Failed() bool {
	return msg.err != nil
}
