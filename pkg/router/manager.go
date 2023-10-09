package router

import (
	interfaces2 "github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
	"net/http"
	"sync"
)

type Manager struct {
	TemplateEngine interfaces2.ITemplateEngine
	templatePath   string
	slugParams     map[string]string
	userContext    sync.Map
	websocket      interfaces2.IWebsocket
}

func NewManager(engine interfaces2.ITemplateEngine) *Manager {
	return &Manager{TemplateEngine: engine}
}

func (m *Manager) SetUserContext(key string, value interface{}) {
	m.userContext.Store(key, value)
}

func (m *Manager) GetUserContext(key string) (any, bool) {
	value, ok := m.userContext.Load(key)
	return value, ok
}

func (m *Manager) SetTemplateEngine(engine interfaces2.ITemplateEngine) {
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
	m.TemplateEngine.SetResponseWriter(w)
	m.TemplateEngine.SetRequest(r)
	err := m.TemplateEngine.Exec()
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

func (m *Manager) GetWebSocket() interfaces2.IWebsocket {
	return m.websocket
}

func (m *Manager) SetWebsocket(websocket interfaces2.IWebsocket) {
	m.websocket = websocket
}
