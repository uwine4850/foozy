package interfaces

import (
	"github.com/uwine4850/foozy/internal/tmlengine"
	"net/http"
)

type IManager interface {
	SetTemplateEngine(engine tmlengine.ITemplateEngine)
	RenderTemplate(w http.ResponseWriter) error
	SetTemplatePath(templatePath string)
	SetContext(data map[string]interface{})
	SetSlugParams(params map[string]string)
	GetSlugParams(key string) (string, bool)
}

type IRouter interface {
	Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	GetMux() *http.ServeMux
	SetTemplateEngine(engine tmlengine.ITemplateEngine)
	SetMiddleware(middleware IMiddleware)
}
