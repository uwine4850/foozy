package form

import "fmt"

type ErrFormConvertFieldNotFound struct {
	Field string
}

func (e ErrFormConvertFieldNotFound) Error() string {
	return fmt.Sprintf("The form field %s was not found.", e.Field)
}
