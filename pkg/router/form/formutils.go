package form

import (
	"fmt"
	"mime/multipart"
	netUrl "net/url"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/pkg/utils/fslice"
)

// FillableFormStruct structure is intended for more convenient access to the structure to be filled in.
// The structure to be filled is passed by pointer, it is filled independently of this structure, so it is up to the
// user to decide whether to use FillableFormStruct.
type FillableFormStruct struct {
	s            interface{}
	defaultValue func(name string) interface{}
}

func NewFillableFormStruct(fillStruct interface{}) *FillableFormStruct {
	if reflect.TypeOf(fillStruct).Kind() != reflect.Pointer {
		panic("the fillStruct parameter must be a reference to a structure")
	}
	return &FillableFormStruct{s: fillStruct}
}

// GetStruct getting a fillable structure.
func (f *FillableFormStruct) GetStruct() interface{} {
	return f.s
}

// SetDefaultValue sets the function that will be executed when the default value is needed.
// If it is not set, the default value will be an empty string.
// func(name string) interface{} - parameter name is the name of the passed key in GetOrDef.
func (f *FillableFormStruct) SetDefaultValue(val func(name string) interface{}) {
	f.defaultValue = val
}

// GetOrDef get slice value or default value if it does not exist.
// name - name of the structure field. Case-sensitive.
// index - index of the structure element.
func (f *FillableFormStruct) GetOrDef(name string, index int) interface{} {
	value := reflect.ValueOf(f.s).Elem()
	fieldValue := value.FieldByName(name)
	if fieldValue.Kind() == reflect.Invalid {
		panic(fmt.Sprintf("field %s not exist", name))
	}
	if fieldValue.IsNil() {
		if f.defaultValue == nil {
			return ""
		}
		return f.defaultValue(name)
	}
	return fieldValue.Index(index).Interface()
}

type FormFile struct {
	Header *multipart.FileHeader
}

// FillStructFromForm A method that fills the structure with data from the form.
// The structure should always be passed as a pointer.
// For correct work it is necessary to specify "form" tag for each field of the structure. For example, `form:<form field name>`.
// Structure fields can be of two types only:
// []FormFile - form files.
// []string - all other data.
// The nilIfNotExist parameter sets the name of form fields that should be nil if they are not found (e.g. useful for checkboxes).
func FillStructFromForm(frm *Form, fillableStruct *FillableFormStruct, nilIfNotExist []string) error {
	fill := fillableStruct.GetStruct()
	if !typeopr.IsPointer(fill) {
		return typeopr.ErrValueNotPointer{Value: "fill"}
	}
	if !typeopr.PtrIsStruct(fill) {
		return typeopr.ErrParameterNotStruct{Param: "fill"}
	}
	orderedForm := FrmValueToOrderedForm(frm)
	t := reflect.TypeOf(fill).Elem()
	v := reflect.ValueOf(fill).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tag := field.Tag.Get("form")
		// Skip if the tag is not a form
		if tag == "" {
			continue
		}
		err := checkFieldType(field)
		if err != nil {
			return err
		}
		orderedFormValue, ok := orderedForm.GetByName(tag)
		formValue := orderedFormValue.Value
		if !ok {
			// Skips loop iteration if the field is not found, but it must be as nil.
			if nilIfNotExist != nil && fslice.SliceContains(nilIfNotExist, tag) {
				continue
			} else {
				return ErrFormConvertFieldNotFound{tag}
			}
		}
		// Set files
		if reflect.DeepEqual(field.Type, reflect.TypeOf([]FormFile{})) && reflect.TypeOf(formValue) == reflect.TypeOf([]FormFile{}) {
			formType, _ := formValue.([]FormFile)
			if !ok {
				return typeopr.ErrConvertType{Type1: reflect.TypeOf(formValue).String(), Type2: "[]FormFile"}
			}
			value.Set(reflect.ValueOf(formType))
		}
		// Set string
		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.String {
			formType, ok := formValue.([]string)
			if !ok {
				return typeopr.ErrConvertType{Type1: reflect.TypeOf(formValue).String(), Type2: "string"}
			}
			value.Set(reflect.ValueOf(formType))
		}
	}
	return nil
}

// checkFieldType checks if the type of field to be filled is correct.
func checkFieldType(field reflect.StructField) error {
	if field.Type.Kind() != reflect.Slice {
		return fmt.Errorf("the %s field should be slice", field.Name)
	}
	if field.Type.Elem() != reflect.TypeOf("") && field.Type.Elem() != reflect.TypeOf(FormFile{}) {
		return fmt.Errorf("the %s field can only be of two types: []string or []FormFile", field.Name)
	}
	return nil
}

type IFormGetEnctypeData interface {
	GetMultipartForm() *multipart.Form
	GetApplicationForm() netUrl.Values
}

// OrderedForm Values can be displayed either by field name or all fields at once.
type OrderedForm struct {
	itemCount int
	names     map[string]int
	values    []OrderedFormValue
}

func NewOrderedForm() *OrderedForm {
	o := &OrderedForm{}
	o.itemCount = 0
	o.names = make(map[string]int)
	return o
}

// Add adds a new form field.
func (f *OrderedForm) Add(name string, value interface{}) {
	f.values = append(f.values, OrderedFormValue{
		Name:  name,
		Value: value,
	})
	f.itemCount++
	f.names[name] = f.itemCount
}

// GetByName getting a field by name.
func (f *OrderedForm) GetByName(name string) (OrderedFormValue, bool) {
	getIndex, ok := f.names[name]
	if !ok {
		return OrderedFormValue{}, ok
	}
	return f.values[getIndex-1], true
}

// GetAll getting all fields.
func (f *OrderedForm) GetAll() []OrderedFormValue {
	return f.values
}

type OrderedFormValue struct {
	Name  string
	Value interface{}
}

// FrmValueToOrderedForm Converts the form to a OrderedForm.
func FrmValueToOrderedForm(frm IFormGetEnctypeData) *OrderedForm {
	orderedForm := NewOrderedForm()
	multipartForm := frm.GetMultipartForm()
	if multipartForm != nil {
		for name, value := range multipartForm.Value {
			orderedForm.Add(name, value)
		}
		for name, value := range multipartForm.File {
			var files []FormFile
			for i := 0; i < len(value); i++ {
				files = append(files, FormFile{Header: value[i]})
			}
			orderedForm.Add(name, files)
		}
	}
	applicationForm := frm.GetApplicationForm()
	if applicationForm != nil {
		for name, value := range applicationForm {
			orderedForm.Add(name, value)
		}
	}
	return orderedForm
}

type ErrArgumentNotPointer struct {
	Name string
}

func (e ErrArgumentNotPointer) Error() string {
	return fmt.Sprintf("argument %s is not a pointer", e.Name)
}

func getStructElemFromFillableStruct(fillableStruct *FillableFormStruct) (reflect.Type, reflect.Value) {
	fillStruct := fillableStruct.GetStruct()
	var fillValue reflect.Value
	var fillType reflect.Type
	// If fillStruct is reflect.Value.
	if reflect.TypeOf(fillStruct) == reflect.TypeOf(&reflect.Value{}) {
		fv := fillStruct.(*reflect.Value)
		fillValue = reflect.Indirect(*fv)
		fillType = fillValue.Type()
	} else {
		// If fillStruct is default struct.
		fillValue = reflect.ValueOf(fillStruct).Elem()
		fillType = reflect.TypeOf(fillStruct).Elem()
	}
	return fillType, fillValue
}

// FieldsNotEmpty checks the specified fields of the structure for emptiness.
// fieldsName - slice with exact names of STRUCTURE fields that should not be empty.
// Optimized to work even if the FillableFormStruct contains a structure with type *reflect.Value.
func FieldsNotEmpty(fillableStruct *FillableFormStruct, fieldsName []string) error {
	fillStruct := fillableStruct.GetStruct()
	if !typeopr.IsPointer(fillStruct) {
		return typeopr.ErrValueNotPointer{Value: "fillableStruct"}
	}
	if !typeopr.PtrIsStruct(fillStruct) {
		return typeopr.ErrParameterNotStruct{Param: "fillableStruct"}
	}

	fillType, fillValue := getStructElemFromFillableStruct(fillableStruct)
	for i := 0; i < len(fieldsName); i++ {
		_, ok := fillType.FieldByName(fieldsName[i])
		if ok {
			val := fillValue.FieldByName(fieldsName[i])
			if val.IsNil() {
				return fmt.Errorf("field %s is empty", fieldsName[i])
			}
			if val.Len() != 0 && val.Type().Elem().Kind() == reflect.String && fslice.AllStringItemsEmpty(val.Interface().([]string)) {
				return fmt.Errorf("field %s is empty", fieldsName[i])
			}
		}
	}
	return nil
}

// FieldsName returns a list of field names of the filled structure.
func FieldsName(fillableStruct *FillableFormStruct, exclude []string) ([]string, error) {
	fillStruct := fillableStruct.GetStruct()
	if !typeopr.IsPointer(fillStruct) {
		return nil, typeopr.ErrValueNotPointer{Value: "fillableStruct"}
	}
	if !typeopr.PtrIsStruct(fillStruct) {
		return nil, typeopr.ErrParameterNotStruct{Param: "fillableStruct"}
	}

	fillType, _ := getStructElemFromFillableStruct(fillableStruct)
	var names []string
	for i := 0; i < fillType.NumField(); i++ {
		field := fillType.Field(i)
		if exclude != nil && fslice.SliceContains(exclude, field.Name) {
			continue
		}
		names = append(names, field.Name)
	}
	return names, nil
}

type ErrExtensionNotMatch struct {
	Field string
}

func (e ErrExtensionNotMatch) Error() string {
	return fmt.Sprintf("The extension of the %s field does not match what is expected.", e.Field)
}

// CheckExtension Check if the file resolution matches the expected one. Can only be used with a structure already
// filled out in the form.
// To work, you need to add an ext tag with the necessary extensions (if there are many, separated by commas).
// For example, ext:".jpg .jpeg .png".
func CheckExtension(fillableStruct *FillableFormStruct) error {
	fillStruct := fillableStruct.GetStruct()
	if !typeopr.IsPointer(fillStruct) {
		return typeopr.ErrValueNotPointer{Value: "fillableStruct"}
	}
	if !typeopr.PtrIsStruct(fillStruct) {
		return typeopr.ErrParameterNotStruct{Param: "fillableStruct"}
	}

	fillType, fillValue := getStructElemFromFillableStruct(fillableStruct)

	for i := 0; i < fillType.NumField(); i++ {
		tag := fillType.Field(i).Tag.Get("ext")
		if tag == "" {
			continue
		}
		field := fillType.Field(i)
		if field.Type.Elem() != reflect.TypeOf(FormFile{}) {
			panic("the ext tag can only be added to fields whose type is form.FormFile")
		}
		extension := strings.Split(tag, " ")
		files := fillValue.Field(i).Interface().([]FormFile)
		for i := 0; i < len(files); i++ {
			if !fslice.SliceContains(extension, filepath.Ext(files[i].Header.Filename)) {
				return ErrExtensionNotMatch{Field: field.Name}
			}
		}
	}
	return nil
}

// FillReflectValueFromForm fills the structure with data from the form.
// It works and does everything the same as the FillStructFromForm function.
// The only difference is that this function accepts data in *reflect.Value format.
func FillReflectValueFromForm(frm *Form, fillValue *reflect.Value, nilIfNotExist []string) error {
	orderedForm := FrmValueToOrderedForm(frm)
	fillType := fillValue.Type()
	for i := 0; i < fillType.NumField(); i++ {
		typeField := fillType.Field(i)
		valueField := fillValue.Field(i)

		tag := typeField.Tag.Get("form")
		if tag == "" {
			continue
		}
		orderedFormValue, ok := orderedForm.GetByName(tag)
		if !ok {
			// Skips loop iteration if the field is not found, but it must be as nil.
			if nilIfNotExist != nil && fslice.SliceContains(nilIfNotExist, tag) {
				continue
			} else {
				return ErrFormConvertFieldNotFound{tag}
			}
		}
		formValue := orderedFormValue.Value

		if reflect.DeepEqual(typeField.Type, reflect.TypeOf([]FormFile{})) && reflect.TypeOf(formValue) == reflect.TypeOf([]FormFile{}) {
			formType, _ := formValue.([]FormFile)
			if !ok {
				return typeopr.ErrConvertType{Type1: reflect.TypeOf(formValue).String(), Type2: "[]FormFile"}
			}
			valueField.Set(reflect.ValueOf(formType))
		}
		// Set string
		if typeField.Type.Kind() == reflect.Slice && typeField.Type.Elem().Kind() == reflect.String {
			formType, ok := formValue.([]string)
			if !ok {
				return typeopr.ErrConvertType{Type1: reflect.TypeOf(formValue).String(), Type2: "string"}
			}
			valueField.Set(reflect.ValueOf(formType))
		}
	}
	return nil
}
