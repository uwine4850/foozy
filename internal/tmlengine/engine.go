package tmlengine

import (
	"github.com/flosch/pongo2"
	"net/http"
)

type TemplateEngine struct {
	path         string
	templateFile *pongo2.Template
	context      map[string]interface{}
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

func (e *TemplateEngine) Exec(w http.ResponseWriter) error {
	err := e.parseFile()
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
	e.context = data
}
