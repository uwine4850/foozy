package interfaces

import (
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces/intrnew"
)

type ITemplateEngine interface {
	intrnew.INewInstance
	SetPath(files string)
	Exec() error
	SetContext(data map[string]interface{})
	GetContext() map[string]interface{}
	SetResponseWriter(w http.ResponseWriter)
	SetRequest(r *http.Request)
}

type IRender interface {
	intrnew.INewInstance
	SetContext(data map[string]interface{})
	GetContext() map[string]interface{}
	SetTemplateEngine(engine ITemplateEngine)
	GetTemplateEngine() ITemplateEngine
	RenderTemplate(w http.ResponseWriter, r *http.Request) error
	SetTemplatePath(templatePath string)
	RenderJson(data interface{}, w http.ResponseWriter) error
}
