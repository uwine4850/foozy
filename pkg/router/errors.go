package router

import "fmt"

type ErrTemplatePathNotSet struct {
}

func (e ErrTemplatePathNotSet) Error() string {
	return "The path to the template is not set."
}

type ErrTemplatePathNotExist struct {
	path string
}

func (e ErrTemplatePathNotExist) Error() string {
	return fmt.Sprintf("The path to the \"%s\" template was not found.", e.path)
}
