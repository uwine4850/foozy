package form

import "fmt"

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

type ErrFormConvertFieldNotFound struct {
	Field string
}

func (e ErrFormConvertFieldNotFound) Error() string {
	return fmt.Sprintf("The form field %s was not found.", e.Field)
}
