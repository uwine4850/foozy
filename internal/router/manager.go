package router

import (
	"github.com/uwine4850/foozy/internal/interfaces"
	"github.com/uwine4850/foozy/internal/utils"
	"net/http"
	"sync"
)

type Manager struct {
	TemplateEngine interfaces.ITemplateEngine
	templatePath   string
	slugParams     map[string]string
	userContext    sync.Map
}

func NewManager(engine interfaces.ITemplateEngine) *Manager {
	return &Manager{TemplateEngine: engine}
}

func (m *Manager) SetUserContext(key string, value interface{}) {
	m.userContext.Store(key, value)
}

func (m *Manager) GetUserContext(key string) (any, bool) {
	value, ok := m.userContext.Load(key)
	return value, ok
}

func (m *Manager) SetTemplateEngine(engine interfaces.ITemplateEngine) {
	m.TemplateEngine = engine
}

// RenderTemplate Rendering a template using a template engine.
func (m *Manager) RenderTemplate(w http.ResponseWriter, r *http.Request) error {
	if m.templatePath == "" {
		return &ErrTemplatePathNotSet{}
	}
	if !utils.PathExist(m.templatePath) {
		return &ErrTemplatePathNotExist{m.templatePath}
	}
	m.TemplateEngine.SetPath(m.templatePath)
	err := m.TemplateEngine.Exec(w, r)
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

// SetSlugParams sets the slug parameters.
func (m *Manager) SetSlugParams(params map[string]string) {
	m.slugParams = params
}

// GetSlugParams returns the parameter by key. If the key is not found returns false.
func (m *Manager) GetSlugParams(key string) (string, bool) {
	res, ok := m.slugParams[key]
	return res, ok
}
