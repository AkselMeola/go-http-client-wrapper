package http_client_wrapper

type Error struct {
	msg      string
	original error
}

func NewError(msg string, original error) *Error {
	return &Error{msg: msg, original: original}
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Unwrap() error {
	return e.original
}
