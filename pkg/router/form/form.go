package form

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/uwine4850/foozy/pkg/router"
	"github.com/uwine4850/foozy/pkg/utils"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Form struct {
	multipartForm      *multipart.Form
	applicationForm    url.Values
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

func (f *Form) GetMultipartForm() *multipart.Form {
	return f.multipartForm
}

func (f *Form) GetApplicationForm() url.Values {
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

func (f *Form) Files(key string) ([]*multipart.FileHeader, bool) {
	fi, ok := f.multipartForm.File[key]
	return fi, ok
}

// ValidateCsrfToken checks the validity of the csrf token. If no errors are detected, the token is valid.
// It is desirable to use this method only after Parse() method.
func (f *Form) ValidateCsrfToken() error {
	csrfToken := f.Value("csrf_token")
	if csrfToken == "" {
		return ErrCsrfTokenNotFound{}
	}
	cookie, err := f.request.Cookie("csrf_token")
	if err != nil {
		return err
	}
	if cookie.Value != csrfToken {
		return ErrCsrfTokenDoesNotMatch{}
	}
	return nil
}

func randomiseTheFileName(pathToDir string, fileName string) string {
	outputFilepath := pathToDir + fileName
	if utils.PathExist(pathToDir + fileName) {
		hash := sha256.Sum256([]byte(fileName))
		hashData := hex.EncodeToString(hash[:])
		ext := filepath.Ext(fileName)
		return randomiseTheFileName(pathToDir, hashData+ext)
	}
	return outputFilepath
}

func SaveFile(w http.ResponseWriter, file multipart.File, fileHeader *multipart.FileHeader, pathToDir string, buildPath *string) bool {
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			router.ServerError(w, err.Error())
		}
	}(file)

	fp := randomiseTheFileName(pathToDir, fileHeader.Filename)
	*buildPath = fp
	dst, err := os.Create(fp)
	if err != nil {
		router.ServerError(w, err.Error())
		return false
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			router.ServerError(w, err.Error())
		}
	}(dst)
	_, err = io.Copy(dst, file)
	if err != nil {
		router.ServerError(w, err.Error())
		return false
	}
	return true
}
