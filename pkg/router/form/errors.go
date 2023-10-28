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

type ErrFormConvertType struct {
	Type1 string
	Type2 string
}

func (e ErrFormConvertType) Error() string {
	return fmt.Sprintf("Form data type conversion error. The %s interface type cannot be converted to %s type.", e.Type1, e.Type2)
}

type ErrFormConvertFieldNotFound struct {
	Field string
}

func (e ErrFormConvertFieldNotFound) Error() string {
	return fmt.Sprintf("The form field %s was not found.", e.Field)
}

type ErrParameterNotPointer struct {
	Param string
}

func (e ErrParameterNotPointer) Error() string {
	return fmt.Sprintf("The %s parameter is not a pointer.", e.Param)
}

type ErrParameterNotStruct struct {
	Param string
}

func (e ErrParameterNotStruct) Error() string {
	return fmt.Sprintf("The %s parameter is not a structure.", e.Param)
}
