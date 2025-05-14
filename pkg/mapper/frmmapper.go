package mapper

import (
	"fmt"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/router/form"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

var formFileType = reflect.TypeOf(form.FormFile{})
var frmRawCache sync.Map
var FC FormConverter

// FillStructFromForm fills the structure with data from the form.
// All standard types and [form.FormFile] for files are supported. Slices
// with all these types are also supported.
// Does not work with nested structures.
//
// Before passing the form to the function, it must be pre-processed using the [Parse] method.
//
// All fields that have the tag `form:"<input_name>"` will be processed.
// Where <input_name> is the name of the key (or HTML input) to the data in the form.
//
// Form field data can exist in three types:
// 1. Just data from the form. For example, strings or numbers.
// 2. Empty data. An empty string "" is passed.
// 3. There is no data. Data is expected, but it is not present in any way.
// There are special tags to handle the first and second cases. Namely:
// __empty__ — executed in the second case, when the string is empty. Tag with empty value is ignored.
// It has several arguments:
//   - `empty:"<some_text>"` - the value that will be passed to the structure field. You can pass, for example,
//     a number, it will be formatted if the field type is int.
//   - `empty:“-err”` - will print the corresponding error if the field is empty.
//
// __nil__ - is executed when it is impossible to find data in the form by the current key, which means that there
// is no data. You can pass the arguments:
//   - `nil:“-skip”` - skips the given field.
//
// If the tag "nil" is not set, there will be an error because of an undiscovered key.
func FillStructFromForm[T any](frm *form.Form, out *T) error {
	of := FrmValueToOrderedForm(frm)

	v := typeopr.GetReflectValue(out)
	raw := LoadSomeRawObjectFromCache(v, &frmRawCache, namelib.TAGS.FORM_MAPPER_NAME)

	for name, f := range *raw.Fields() {
		fieldValue := v.FieldByName(f.Name)
		orderedFormValues, ok := of.GetByName(name)
		if !ok {
			switch f.Tag.Get(namelib.TAGS.FORM_MAPPER_NIL) {
			case NIL_SKIP:
				continue
			default:
				return form.ErrFormConvertFieldNotFound{Field: name}
			}
		}
		if err := FC.handleItem(&orderedFormValues, &fieldValue, f.Name, f.Tag.Get(namelib.TAGS.FORM_MAPPER_EMPTY)); err != nil {
			return err
		}
	}
	if err := CheckExtension(out); err != nil {
		return err
	}
	return nil
}

// CheckExtension Check if the file resolution matches the expected one. Can only be used with a structure already
// filled out in the form.
// To work, you need to add an ext tag with the necessary extensions (if there are many, separated by commas).
// For example, ext:".jpg .jpeg .png".
func CheckExtension[T any](filledStruct *T) error {
	v := typeopr.GetReflectValue(filledStruct)
	raw := LoadSomeRawObjectFromCache(v, &frmRawCache, namelib.TAGS.FORM_MAPPER_NAME)

	for _, f := range *raw.Fields() {
		fieldValue := v.FieldByName(f.Name)
		extensionsString := f.Tag.Get(namelib.TAGS.FORM_MAPPER_EXTENSION)
		if extensionsString == "" {
			continue
		}
		extensionSlice := strings.Split(extensionsString, " ")
		switch f.Type.Kind() {
		case reflect.Slice:
			if f.Type.Elem() != formFileType {
				panic("the ext tag can only be added to fields whose type is form.FormFile")
			}
			files := fieldValue.Interface().([]form.FormFile)
			for i := 0; i < len(files); i++ {
				if !slices.Contains(extensionSlice, filepath.Ext(files[i].Header.Filename)) {
					return ErrExtensionNotMatch{Field: f.Name}
				}
			}
		case reflect.Struct:
			if f.Type != formFileType {
				panic("the ext tag can only be added to fields whose type is form.FormFile")
			}
			fmt.Println(fieldValue)
			file := fieldValue.Interface().(form.FormFile)
			if !slices.Contains(extensionSlice, filepath.Ext(file.Header.Filename)) {
				return ErrExtensionNotMatch{Field: f.Name}
			}
		}
	}
	return nil
}

type ErrExtensionNotMatch struct {
	Field string
}

func (e ErrExtensionNotMatch) Error() string {
	return fmt.Sprintf("The extension of the %s field does not match what is expected.", e.Field)
}
