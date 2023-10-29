package dbutils

import "fmt"

type ErrDbResFieldNotFound struct {
	Field string
}

func (e ErrDbResFieldNotFound) Error() string {
	return fmt.Sprintf("The database output result does not contain the %s field.", e.Field)
}
