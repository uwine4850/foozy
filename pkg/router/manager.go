package router

import (
	"encoding/json"
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

// SetUserContext sets the user context.
// This context is used only as a means of passing information between handlers.
func (m *Manager) SetUserContext(key string, value interface{}) {
	m.userContext.Store(key, value)
}

// GetUserContext getting the user context.
func (m *Manager) GetUserContext(key string) (any, bool) {
	value, ok := m.userContext.Load(key)
	return value, ok
}

// DelUserContext deletes a user context by key.
func (m *Manager) DelUserContext(key string) {
	m.userContext.Delete(key)
}

// SetTemplateEngine set the template engine interface.
// Optional method if the template engine is already installed.
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

// GetWebSocket get an instance for the websocket connection.
// Works only in the "Ws" handler.
func (m *Manager) GetWebSocket() interfaces2.IWebsocket {
	return m.websocket
}

// SetWebsocket sets the websocket interface.
func (m *Manager) SetWebsocket(websocket interfaces2.IWebsocket) {
	m.websocket = websocket
}

// RenderJson displays data in json format on the page.
func (m *Manager) RenderJson(data interface{}, w http.ResponseWriter) error {
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
