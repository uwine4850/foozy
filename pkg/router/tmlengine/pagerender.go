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

func NewRender() (interfaces.IRender, error) {
	render := Render{}
	newRender, err := render.New()
	if err != nil {
		return nil, err
	}
	return newRender.(interfaces.IRender), nil
}

func (rn *Render) New() (interface{}, error) {
	var engine interfaces.ITemplateEngine
	if rn.TemplateEngine != nil {
		engine = rn.TemplateEngine
	} else {
		engine = &TemplateEngine{}
	}
	if engine, err := engine.New(); err != nil {
		return nil, err
	} else {
		return &Render{TemplateEngine: engine.(interfaces.ITemplateEngine)}, nil
	}
}

// SetContext Setting variables for html template.
func (rn *Render) SetContext(data map[string]interface{}) {
	rn.TemplateEngine.SetContext(data)
}

func (rn *Render) GetContext() map[string]interface{} {
	return rn.TemplateEngine.GetContext()
}

// SetTemplateEngine set the template engine interface.
// Optional method if the template engine is already installed.
func (rn *Render) SetTemplateEngine(engine interfaces.ITemplateEngine) {
	rn.TemplateEngine = engine
}

func (rn *Render) GetTemplateEngine() interfaces.ITemplateEngine {
	return rn.TemplateEngine
}

// RenderTemplate Rendering a template using a template engine.
func (rn *Render) RenderTemplate(w http.ResponseWriter, r *http.Request) error {
	if rn.templatePath == "" {
		return ErrTemplatePathNotSet{}
	}
	if !utils.PathExist(rn.templatePath) {
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

// SetTemplatePath Setting the path to the template that the templating engine renders.
func (rn *Render) SetTemplatePath(templatePath string) {
	rn.templatePath = templatePath
}

// RenderJson displays data in json format on the page.
func (rn *Render) RenderJson(data interface{}, w http.ResponseWriter) error {
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

// CreateAndSetNewRenderInstance creates and sets a new Render instance into the manager.
func CreateAndSetNewRenderInstance(manager interfaces.IManager) error {
	render := manager.Render()

	var newRender interfaces.IRender
	err := utils.CreateNewInstance(render, &newRender)
	if err != nil {
		return err
	}
	tmplEngine := render.GetTemplateEngine()

	var newTmplEngine interfaces.ITemplateEngine
	err = utils.CreateNewInstance(tmplEngine, &newTmplEngine)
	if err != nil {
		return err
	}

	newRender.SetTemplateEngine(newTmplEngine)
	manager.SetRender(newRender)
	return nil
}
