package manager

import (
	"encoding/json"
	"net/http"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
)

type ManagerRender struct {
	TemplateEngine interfaces.ITemplateEngine
	templatePath   string
}

func NewManagerRender(engine interfaces.ITemplateEngine) *ManagerRender {
	return &ManagerRender{TemplateEngine: engine}
}

// SetContext Setting variables for html template.
func (m *ManagerRender) SetContext(data map[string]interface{}) {
	m.TemplateEngine.SetContext(data)
}

// SetTemplateEngine set the template engine interface.
// Optional method if the template engine is already installed.
func (m *ManagerRender) SetTemplateEngine(engine interfaces.ITemplateEngine) {
	m.TemplateEngine = engine
}

// RenderTemplate Rendering a template using a template engine.
func (m *ManagerRender) RenderTemplate(w http.ResponseWriter, r *http.Request) error {
	if m.templatePath == "" {
		// return &router.ErrTemplatePathNotSet{}
		panic("AA")
	}
	if !utils.PathExist(m.templatePath) {
		// return &router.ErrTemplatePathNotExist{Path: m.templatePath}
		panic("!!!!")
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
func (m *ManagerRender) SetTemplatePath(templatePath string) {
	m.templatePath = templatePath
}

// RenderJson displays data in json format on the page.
func (m *ManagerRender) RenderJson(data interface{}, w http.ResponseWriter) error {
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
