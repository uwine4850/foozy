package interfaces

import (
	"net/http"
)

type IManager interface {
	SetTemplateEngine(engine ITemplateEngine)
	RenderTemplate(w http.ResponseWriter) error
	SetTemplatePath(templatePath string)
	SetContext(data map[string]interface{})
	SetSlugParams(params map[string]string)
	GetSlugParams(key string) (string, bool)
	SetUserContext(key string, value interface{})
	GetUserContext(key string) (any, bool)
}

type IRouter interface {
	Get(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	Post(pattern string, fn func(w http.ResponseWriter, r *http.Request, manager IManager))
	GetMux() *http.ServeMux
	SetTemplateEngine(engine ITemplateEngine)
	SetMiddleware(middleware IMiddleware)
}
