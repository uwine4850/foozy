package form

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	netUrl "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/uwine4850/foozy/pkg/interfaces"
)

// FormFile representation of the form file as an object.
// All data of the file is stored here.
type FormFile struct {
	Header *multipart.FileHeader
}

type Form struct {
	multipartForm      *multipart.Form
	applicationForm    netUrl.Values
	multipartMaxMemory int64
	request            *http.Request
}

func NewForm(request *http.Request) *Form {
	return &Form{multipartMaxMemory: 32 << 20, request: request}
}

// Parse parsing a form method. After that you can access fields by name and the form in general.
func (f *Form) Parse() error {
	enctype := strings.Split(f.request.Header.Get("Content-Type"), ";")[0]
	switch enctype {
	case "application/x-www-form-urlencoded":
		err := f.request.ParseForm()
		if err != nil {
			return err
		}
		f.applicationForm = f.request.PostForm
	case "multipart/form-data":
		err := f.request.ParseMultipartForm(f.multipartMaxMemory)
		if err != nil {
			return err
		}
		f.multipartForm = f.request.MultipartForm
	}
	return nil
}

// GetMultipartForm getting multipart/form-data form data.
func (f *Form) GetMultipartForm() *multipart.Form {
	return f.multipartForm
}

// GetApplicationForm getting application/x-www-form-urlencoded form data.
func (f *Form) GetApplicationForm() netUrl.Values {
	return f.applicationForm
}

// Value getting a simple value of the form.
func (f *Form) Value(key string) string {
	val := f.request.PostFormValue(key)
	return template.HTMLEscapeString(val)
}

// File retrieving a file from a form. Multipart/form-data only.
func (f *Form) File(key string) (multipart.File, *multipart.FileHeader, error) {
	return f.request.FormFile(key)
}

// Files getting multiple files from multipart type input.
func (f *Form) Files(key string) ([]*multipart.FileHeader, bool) {
	fi, ok := f.multipartForm.File[key]
	return fi, ok
}

// randomiseTheFileName If the file name already exists, randomises it and returns the new file path.
func randomiseFileName(fileName string) string {
	ext := filepath.Ext(fileName)
	base := fileName[:len(fileName)-len(ext)]
	uniqueID := uuid.New().String()
	return base + "-" + uniqueID + ext
}

// SaveFile Saves the file in the specified directory.
// If the file name is already found, uses the randomiseTheFileName function to randomise the file name.
func SaveFile(fileHeader *multipart.FileHeader, pathToDir string, buildPath *string, manager interfaces.Manager) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	fn := randomiseFileName(fileHeader.Filename)
	newFilePath := filepath.Join(pathToDir, fn)
	if buildPath != nil {
		*buildPath = newFilePath
	}
	dst, err := os.Create(newFilePath)
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, file)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	// Close dst.
	if err := dst.Close(); err != nil {
		return err
	}
	return nil
}

// ReplaceFile changes the specified file to a new file.
func ReplaceFile(pathToFile string, fileHeader *multipart.FileHeader, pathToDir string, buildPath *string, manager interfaces.Manager) error {
	err := os.Remove(pathToFile)
	if err != nil {
		return err
	}
	err = SaveFile(fileHeader, pathToDir, buildPath, manager)
	if err != nil {
		return err
	}
	return nil
}

// SendApplicationForm sends a form of type application/x-www-form-urlencoded to the specified address.
func SendApplicationForm(url string, values map[string][]string) (*http.Response, error) {
	formData := netUrl.Values{}
	for name, value := range values {
		for i := 0; i < len(value); i++ {
			formData.Add(name, value[i])
		}
	}
	response, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, err
	}
	return response, nil
}

// SendMultipartForm sends a form of type multipart/form-data to the specified address.
// The files argument accepts form field names and a slice with file paths.
func SendMultipartForm(url string, values map[string][]string, files map[string][]string) (*http.Response, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for name, value := range files {
		for i := 0; i < len(value); i++ {
			if value[i] == "" {
				fileWriter, err := writer.CreateFormFile(name, "")
				if err != nil {
					return nil, err
				}
				if _, err := io.Copy(fileWriter, bytes.NewReader(nil)); err != nil {
					return nil, err
				}
				continue
			}
			file, err := os.Open(value[i])
			if err != nil {
				return nil, err
			}
			defer file.Close()
			fileWriter, err := writer.CreateFormFile(name, value[i])
			if err != nil {
				return nil, err
			}
			if _, err := io.Copy(fileWriter, file); err != nil {
				return nil, err
			}
		}
	}
	for name, value := range values {
		for i := 0; i < len(value); i++ {
			writer.WriteField(name, value[i])
		}
	}
	writer.Close()
	response, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type ErrFormConvertFieldNotFound struct {
	Field string
}

func (e ErrFormConvertFieldNotFound) Error() string {
	return fmt.Sprintf("The form field [%s] was not found.", e.Field)
}
