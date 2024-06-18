package interfaces

import "net/http"

type ITemplateEngine interface {
	INewInstance
	SetPath(files string)
	Exec() error
	SetContext(data map[string]interface{})
	GetContext() map[string]interface{}
	SetResponseWriter(w http.ResponseWriter)
	SetRequest(r *http.Request)
}

type IRender interface {
	INewInstance
	SetContext(data map[string]interface{})
	GetContext() map[string]interface{}
	SetTemplateEngine(engine ITemplateEngine)
	GetTemplateEngine() ITemplateEngine
	RenderTemplate(w http.ResponseWriter, r *http.Request) error
	SetTemplatePath(templatePath string)
	RenderJson(data interface{}, w http.ResponseWriter) error
}
