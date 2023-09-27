package router

import (
	"github.com/uwine4850/foozy/internal/tmlengine"
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

func (m *Manager) RenderTemplate(w http.ResponseWriter) error {
	m.TemplateEngine.SetPath(m.templatePath)
	err := m.TemplateEngine.Exec(w)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) SetTemplatePath(templatePath string) {
	m.templatePath = templatePath
}

func (m *Manager) SetContext(data map[string]interface{}) {
	m.TemplateEngine.SetContext(data)
}
