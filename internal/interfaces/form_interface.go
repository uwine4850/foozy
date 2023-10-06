package interfaces

import (
	"mime/multipart"
	"net/url"
)

type IForm interface {
	Parse() error
	GetMultipartForm() *multipart.Form
	GetApplicationForm() url.Values
	Value(key string) string
	File(key string) (multipart.File, *multipart.FileHeader, error)
	ValidateCsrfToken() error
}
