package manager

import (
	"fmt"

	"github.com/uwine4850/foozy/pkg/interfaces"
)

type Manager struct {
	ManagerRender
	ManagerWebsocket
	ManagerConf
	ManagerData *ManagerData
}

func (m *Manager) Render() interfaces.IManagerRender {
	return &m.ManagerRender
}

func (m *Manager) WS() interfaces.IManagerWebsocket {
	return &m.ManagerWebsocket
}

func (m *Manager) OneTimeData() interfaces.IManagerOneTimeData {
	return m.ManagerData
}

func (m *Manager) Config() interfaces.IManagerConfig {
	return &m.ManagerConf
}

func (m *Manager) SetOneTimeData(manager interfaces.IManagerOneTimeData) {
	if data, ok := manager.(*ManagerData); ok {
		m.ManagerData = data
	} else {
		fmt.Println("Invalid manager type")
	}
}

func NewManager(engine interfaces.ITemplateEngine) *Manager {
	return &Manager{
		ManagerConf:      *NewManagerConf(),
		ManagerData:      NewManagerData(),
		ManagerRender:    *NewManagerRender(engine),
		ManagerWebsocket: *NewManagerWebsocket(),
	}
}
