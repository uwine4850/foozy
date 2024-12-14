package tmlengine

import (
	"encoding/json"
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils/fstring"
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
	var err error
	if rn.TemplateEngine != nil {
		_engine, err := rn.TemplateEngine.New()
		if err != nil {
			return nil, err
		}
		engine = _engine.(interfaces.ITemplateEngine)
	} else {
		engine, err = NewTemplateEngine()
		if err != nil {
			return nil, err
		}
	}
	return &Render{TemplateEngine: engine}, nil
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
func (rn *Render) RenderTemplate(w http.ResponseWriter, r *http.Request, managerConfig interfaces.IManagerConfig) error {
	if rn.templatePath == "" {
		return ErrTemplatePathNotSet{}
	}
	if !fstring.PathExist(rn.templatePath) {
		return ErrTemplatePathNotExist{Path: rn.templatePath}
	}
	rn.TemplateEngine.SetPath(rn.templatePath)
	rn.TemplateEngine.SetResponseWriter(w)
	rn.TemplateEngine.SetRequest(r)
	err := rn.TemplateEngine.Exec(managerConfig)
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
func CreateNewRenderInstance(manager interfaces.IManager) (interfaces.IRender, error) {
	render := manager.Render()
	newRender, err := render.New()
	if err != nil {
		return nil, err
	}
	return newRender.(interfaces.IRender), nil
}
