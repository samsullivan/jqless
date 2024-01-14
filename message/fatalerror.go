package message

type FatalError struct {
	err error
}

func NewFatalError(err error) FatalError {
	return FatalError{err: err}
}

func (msg FatalError) Error() string {
	return msg.err.Error()
}
