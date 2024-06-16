package tmlengine

import (
	"encoding/json"
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
)

type Render struct {
	TemplateEngine interfaces.ITemplateEngine
	templatePath   string
}

func NewRender() (*Render, error) {
	if engine, err := NewTemplateEngine(); err != nil {
		return nil, err
	} else {
		return &Render{TemplateEngine: engine}, nil
	}
}

// SetContext Setting variables for html template.
func (m *Render) SetContext(data map[string]interface{}) {
	m.TemplateEngine.SetContext(data)
}

// SetTemplateEngine set the template engine interface.
// Optional method if the template engine is already installed.
func (m *Render) SetTemplateEngine(engine interfaces.ITemplateEngine) {
	m.TemplateEngine = engine
}

// RenderTemplate Rendering a template using a template engine.
func (m *Render) RenderTemplate(w http.ResponseWriter, r *http.Request) error {
	if m.templatePath == "" {
		return ErrTemplatePathNotSet{}
	}
	if !utils.PathExist(m.templatePath) {
		return ErrTemplatePathNotExist{Path: m.templatePath}
	}
	m.TemplateEngine.SetPath(m.templatePath)
	m.TemplateEngine.SetResponseWriter(w)
	m.TemplateEngine.SetRequest(r)
	err := m.TemplateEngine.Exec()
	if err != nil {
		return err
	}
	return nil
}

// SetTemplatePath Setting the path to the template that the templating engine renders.
func (m *Render) SetTemplatePath(templatePath string) {
	m.templatePath = templatePath
}

// RenderJson displays data in json format on the page.
func (m *Render) RenderJson(data interface{}, w http.ResponseWriter) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(marshal)
	if err != nil {
		return err
	}
	return nil
}
