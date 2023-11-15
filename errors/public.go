package errors

// Public is a wrapper for errors that should be exposed to the user.
// The error can also be unwrapped using the Unwrap method.
func Public(err error, msg string) error {
	return publicError{err: err, msg: msg}
}

type publicError struct {
	err error
	msg string
}

func (pe publicError) Error() string {
	return pe.err.Error()
}

func (pe publicError) Public() string {
	return pe.msg
}

func (pe publicError) Unwrap() error {
	return pe.err
}
