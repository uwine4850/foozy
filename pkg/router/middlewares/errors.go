package middlewares

import (
	"fmt"
	"strconv"
)

type ErrIdAlreadyExist struct {
	id int
}

func (e ErrIdAlreadyExist) Error() string {
	return fmt.Sprintf("Middleware with id %s already exists.", strconv.Itoa(e.id))
}
