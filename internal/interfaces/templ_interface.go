package interfaces

import "net/http"

type ITemplateEngine interface {
	SetPath(files string)
	Exec(w http.ResponseWriter) error
	SetContext(data map[string]interface{})
}
