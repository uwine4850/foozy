package router

import (
	"github.com/uwine4850/foozy/internal/tmlengine"
	"github.com/uwine4850/foozy/internal/utils"
	"net/http"
)

type IManager interface {
	SetTemplateEngine(engine tmlengine.ITemplateEngine)
	RenderTemplate(w http.ResponseWriter) error
	SetTemplatePath(templatePath string)
	SetContext(data map[string]interface{})
}

type Manager struct {
	TemplateEngine tmlengine.ITemplateEngine
	templatePath   string
}

func NewManager() *Manager {
	return &Manager{TemplateEngine: &tmlengine.TemplateEngine{}}
}

func (m *Manager) SetTemplateEngine(engine tmlengine.ITemplateEngine) {
	m.TemplateEngine = engine
}

// RenderTemplate Rendering a template using a template engine.
func (m *Manager) RenderTemplate(w http.ResponseWriter) error {
	if m.templatePath == "" {
		return &ErrTemplatePathNotSet{}
	}
	if !utils.PathExist(m.templatePath) {
		return &ErrTemplatePathNotExist{m.templatePath}
	}
	m.TemplateEngine.SetPath(m.templatePath)
	err := m.TemplateEngine.Exec(w)
	if err != nil {
		return err
	}
	return nil
}

// SetTemplatePath Setting the path to the template that the templating engine renders.
func (m *Manager) SetTemplatePath(templatePath string) {
	m.templatePath = templatePath
}

// SetContext Setting variables for html template.
func (m *Manager) SetContext(data map[string]interface{}) {
	m.TemplateEngine.SetContext(data)
}
