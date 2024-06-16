package manager

import (
	"github.com/uwine4850/foozy/pkg/ferrors"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
)

type Manager struct {
	managerWebsocket interfaces.IManagerWebsocket
	managerConf      interfaces.IManagerConfig
	managerData      interfaces.IManagerOneTimeData
	render           interfaces.IRender
}

func (m *Manager) Render() interfaces.IRender {
	return m.render
}

func (m *Manager) SetRender(render interfaces.IRender) {
	if !utils.IsPointer(render) {
		panic(ferrors.ErrValueNotPointer{ValueName: "render"})
	}
	m.render = render
}

func (m *Manager) SetWS(ws interfaces.IManagerWebsocket) {
	if !utils.IsPointer(ws) {
		panic(ferrors.ErrValueNotPointer{ValueName: "ws"})
	}
	m.managerWebsocket = ws
}

func (m *Manager) WS() interfaces.IManagerWebsocket {
	return m.managerWebsocket
}

func (m *Manager) OneTimeData() interfaces.IManagerOneTimeData {
	return m.managerData
}

func (m *Manager) SetConfig(cnf interfaces.IManagerConfig) {
	if !utils.IsPointer(cnf) {
		panic(ferrors.ErrValueNotPointer{ValueName: "cnf"})
	}
	m.managerConf = cnf
}

func (m *Manager) Config() interfaces.IManagerConfig {
	return m.managerConf
}

func (m *Manager) SetOneTimeData(manager interfaces.IManagerOneTimeData) {
	if !utils.IsPointer(manager) {
		panic(ferrors.ErrValueNotPointer{ValueName: "manager"})
	}
	m.managerData = manager
}

func NewManager(render interfaces.IRender) *Manager {
	return &Manager{
		managerConf:      NewManagerConf(),
		managerData:      NewManagerData(),
		render:           render,
		managerWebsocket: NewManagerWebsocket(),
	}
}
