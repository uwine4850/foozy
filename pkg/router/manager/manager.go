package manager

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

type Manager struct {
	managerConf interfaces.IManagerConfig
	managerData interfaces.IManagerOneTimeData
	render      interfaces.IRender
}

func (m *Manager) Render() interfaces.IRender {
	return m.render
}

func (m *Manager) SetRender(render interfaces.IRender) {
	if !typeopr.IsPointer(render) {
		panic(typeopr.ErrValueNotPointer{Value: "render"})
	}
	m.render = render
}

func (m *Manager) SetOneTimeData(manager interfaces.IManagerOneTimeData) {
	if !typeopr.IsPointer(manager) {
		panic(typeopr.ErrValueNotPointer{Value: "manager"})
	}
	m.managerData = manager
}

func (m *Manager) OneTimeData() interfaces.IManagerOneTimeData {
	return m.managerData
}

func (m *Manager) SetConfig(cnf interfaces.IManagerConfig) {
	if !typeopr.IsPointer(cnf) {
		panic(typeopr.ErrValueNotPointer{Value: "cnf"})
	}
	m.managerConf = cnf
}

func (m *Manager) Config() interfaces.IManagerConfig {
	return m.managerConf
}

func NewManager(render interfaces.IRender) *Manager {
	return &Manager{
		managerConf: NewManagerConf(),
		managerData: NewManagerData(),
		render:      render,
	}
}
