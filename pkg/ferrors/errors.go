package ferrors

import "fmt"

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

type ErrConvertType struct {
	Type1 string
	Type2 string
}

func (e ErrConvertType) Error() string {
	return fmt.Sprintf("Data type conversion error. The %s interface type cannot be converted to %s type.", e.Type1, e.Type2)
}

type ErrUnsupportedTypeConvert struct {
	Type string
}

func (e ErrUnsupportedTypeConvert) Error() string {
	return fmt.Sprintf("Unsupported type: %s", e.Type)
}
