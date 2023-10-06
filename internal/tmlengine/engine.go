package tmlengine

import (
	"errors"
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/uwine4850/foozy/internal/utils"
	"net/http"
)

type TemplateEngine struct {
	path         string
	templateFile *pongo2.Template
	context      map[string]interface{}
}

func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{context: make(map[string]interface{})}
}

func (e *TemplateEngine) SetPath(path string) {
	e.path = path
}

func (e *TemplateEngine) parseFile() error {
	file, err := pongo2.FromFile(e.path)
	if err != nil {
		return err
	}
	e.templateFile = file
	return nil
}

func (e *TemplateEngine) Exec(w http.ResponseWriter, r *http.Request) error {
	err := e.parseFile()
	if err != nil {
		return err
	}
	err = e.setCsrfVariable(r)
	if err != nil {
		return err
	}
	execute, err := e.templateFile.Execute(e.context)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(execute))
	if err != nil {
		return err
	}
	return nil
}

func (e *TemplateEngine) SetContext(data map[string]interface{}) {
	utils.MergeMap(&e.context, data)
}

// setCsrfVariable sets the csrf token as a variable for the templating context.
func (e *TemplateEngine) setCsrfVariable(r *http.Request) error {
	token, err := r.Cookie("csrf_token")
	data := make(map[string]interface{})
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return err
	}
	if errors.Is(err, http.ErrNoCookie) {
		data["csrf_token"] = ""
		e.SetContext(data)
		return nil
	}
	data["csrf_token"] = fmt.Sprintf("<input name=\"csrf_token\" type=\"hidden\" value=\"%s\">", token.Value)
	e.SetContext(data)
	return nil
}
