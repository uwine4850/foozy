package auth

import "fmt"

type ErrUserAlreadyExist struct {
	Username string
}

func (e ErrUserAlreadyExist) Error() string {
	return fmt.Sprintf("User %s already exist.", e.Username)
}

type ErrShortPassword struct {
}

func (e ErrShortPassword) Error() string {
	return "The password must be equal to or longer than 6 characters."
}

type ErrUserNotExist struct {
	Username string
}

func (e ErrUserNotExist) Error() string {
	return fmt.Sprintf("User %s not exist.", e.Username)
}

type ErrPasswordsDontMatch struct {
}

func (e ErrPasswordsDontMatch) Error() string {
	return "The passwords don't match."
}
