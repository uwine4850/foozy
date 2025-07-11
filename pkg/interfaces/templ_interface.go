package interfaces

import (
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces/itypeopr"
)

type TemplateEngine interface {
	itypeopr.NewInstance
	SetPath(files string)
	Exec() error
	SetContext(data map[string]interface{})
	GetContext() map[string]interface{}
	SetResponseWriter(w http.ResponseWriter)
	SetRequest(r *http.Request)
}

type Render interface {
	itypeopr.NewInstance
	SetContext(data map[string]interface{})
	GetContext() map[string]interface{}
	SetTemplateEngine(engine TemplateEngine)
	GetTemplateEngine() TemplateEngine
	RenderTemplate(w http.ResponseWriter, r *http.Request) error
	SetTemplatePath(templatePath string)
	RenderJson(data interface{}, w http.ResponseWriter) error
}
