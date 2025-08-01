package tmlengine

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/flosch/pongo2"
	"github.com/uwine4850/foozy/pkg/debug"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/namelib"
	"github.com/uwine4850/foozy/pkg/utils/fmap"
)

type TemplateEngine struct {
	path         string
	templateFile *pongo2.Template
	context      map[string]interface{}
	writer       http.ResponseWriter
	request      *http.Request
	mu           sync.Mutex
}

func NewTemplateEngine() interfaces.TemplateEngine {
	RegisterMultipleGlobalFilter(BuiltinFilters)
	return &TemplateEngine{context: make(map[string]interface{})}
}

func (e *TemplateEngine) New() (interface{}, error) {
	RegisterMultipleGlobalFilter(BuiltinFilters)
	return &TemplateEngine{context: make(map[string]interface{})}, nil
}

// SetPath sets the path to the template.
func (e *TemplateEngine) SetPath(path string) {
	e.path = path
}

// processingFile processes the template file.
func (e *TemplateEngine) processingFile() error {
	file, err := pongo2.FromFile(e.path)
	if err != nil {
		return err
	}
	e.templateFile = file
	return nil
}

// Exec does all the necessary processing for the template and shows the HTML code on the page.
func (e *TemplateEngine) Exec() error {
	debug.RequestLogginIfEnable(debug.P_TEMPLATE_ENGINE, "exec template engine...")
	debug.RequestLogginIfEnable(debug.P_TEMPLATE_ENGINE, "processing html file")
	err := e.processingFile()
	if err != nil {
		return err
	}
	debug.RequestLogginIfEnable(debug.P_TEMPLATE_ENGINE, "set CSRF token")
	err = e.setCsrfVariable(e.request)
	if err != nil {
		return err
	}
	debug.RequestLogginIfEnable(debug.P_TEMPLATE_ENGINE, "execute template")
	execute, err := e.templateFile.Execute(e.context)
	if err != nil {
		return err
	}
	debug.RequestLogginIfEnable(debug.P_TEMPLATE_ENGINE, "write template")
	_, err = e.writer.Write([]byte(execute))
	if err != nil {
		return err
	}
	return nil
}

// SetContext sets the variables for the template.
func (e *TemplateEngine) SetContext(data map[string]interface{}) {
	fmap.MergeMapSync(&e.mu, &e.context, data)
}

func (e *TemplateEngine) GetContext() map[string]interface{} {
	return e.context
}

func (e *TemplateEngine) SetResponseWriter(w http.ResponseWriter) {
	e.writer = w
}

func (e *TemplateEngine) SetRequest(r *http.Request) {
	e.request = r
}

// setCsrfVariable sets the csrf token as a variable for the templating context.
func (e *TemplateEngine) setCsrfVariable(r *http.Request) error {
	token, err := r.Cookie(namelib.ROUTER.COOKIE_CSRF_TOKEN)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return err
	}
	data := make(map[string]interface{})
	if errors.Is(err, http.ErrNoCookie) {
		e.SetContext(data)
		return nil
	}
	data[namelib.ROUTER.COOKIE_CSRF_TOKEN] = fmt.Sprintf("<input name=\"%s\" type=\"hidden\" value=\"%s\">", namelib.ROUTER.COOKIE_CSRF_TOKEN, token.Value)
	e.SetContext(data)
	return nil
}
