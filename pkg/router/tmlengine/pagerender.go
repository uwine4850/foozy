package tmlengine

import (
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils/fpath"
)

type Render struct {
	TemplateEngine interfaces.TemplateEngine
	templatePath   string
}

func NewRender() (interfaces.Render, error) {
	render := Render{}
	newRender, err := render.New()
	if err != nil {
		return nil, err
	}
	return newRender.(interfaces.Render), nil
}

func (rn *Render) New() (interface{}, error) {
	var engine interfaces.TemplateEngine
	if rn.TemplateEngine != nil {
		_engine, err := rn.TemplateEngine.New()
		if err != nil {
			return nil, err
		}
		engine = _engine.(interfaces.TemplateEngine)
	} else {
		engine = NewTemplateEngine()
	}
	return &Render{TemplateEngine: engine}, nil
}

// SetContext setting variables for html template.
func (rn *Render) SetContext(data map[string]interface{}) {
	rn.TemplateEngine.SetContext(data)
}

func (rn *Render) GetContext() map[string]interface{} {
	return rn.TemplateEngine.GetContext()
}

// SetTemplateEngine set the template engine interface.
// Optional method if the template engine is already installed.
func (rn *Render) SetTemplateEngine(engine interfaces.TemplateEngine) {
	rn.TemplateEngine = engine
}

func (rn *Render) GetTemplateEngine() interfaces.TemplateEngine {
	return rn.TemplateEngine
}

// RenderTemplate Rendering a template using a template engine.
func (rn *Render) RenderTemplate(w http.ResponseWriter, r *http.Request) error {
	if rn.templatePath == "" {
		return ErrTemplatePathNotSet{}
	}
	if !fpath.PathExist(rn.templatePath) {
		return ErrTemplatePathNotExist{Path: rn.templatePath}
	}
	rn.TemplateEngine.SetPath(rn.templatePath)
	rn.TemplateEngine.SetResponseWriter(w)
	rn.TemplateEngine.SetRequest(r)
	err := rn.TemplateEngine.Exec()
	if err != nil {
		return err
	}
	return nil
}

// SetTemplatePath setting the path to the template that the templating engine renders.
func (rn *Render) SetTemplatePath(templatePath string) {
	rn.templatePath = templatePath
}
