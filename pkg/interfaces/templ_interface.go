package interfaces

import "net/http"

type ITemplateEngine interface {
	SetPath(files string)
	Exec() error
	SetContext(data map[string]interface{})
	SetResponseWriter(w http.ResponseWriter)
	SetRequest(r *http.Request)
}
