package manager

import (
	interfaces2 "github.com/uwine4850/foozy/pkg/interfaces"
)

type Manager struct {
	ManagerConf
	ManagerData
	ManagerRender
	ManagerWebsocket
}

func NewManager(engine interfaces2.ITemplateEngine) *Manager {
	return &Manager{
		ManagerConf:      *NewManagerConf(),
		ManagerData:      *NewManagerData(),
		ManagerRender:    *NewManagerRender(engine),
		ManagerWebsocket: *NewManagerWebsocket(),
	}
}
