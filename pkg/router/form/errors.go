package form

type ErrCsrfTokenNotFound struct {
}

func (receiver ErrCsrfTokenNotFound) Error() string {
	return "Csrf token not found."
}

type ErrCsrfTokenDoesNotMatch struct {
}

func (receiver ErrCsrfTokenDoesNotMatch) Error() string {
	return "Csrf token does not match. The validity time may have expired."
}
