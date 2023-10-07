package form

import (
	"html/template"
	"mime/multipart"
	"net/http"
	"net/url"
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
