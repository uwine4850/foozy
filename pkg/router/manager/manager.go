package manager

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/secure"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

type Manager struct {
	managerData interfaces.IManagerOneTimeData
	render      interfaces.IRender
	key         interfaces.IKey
}

func (m *Manager) New() (interface{}, error) {
	if m.key != nil {
		return &Manager{key: m.key}, nil
	} else {
		return &Manager{}, nil
	}
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

func (m *Manager) Key() interfaces.IKey {
	return m.key
}

func NewManager(render interfaces.IRender) *Manager {
	return &Manager{
		managerData: NewManagerData(),
		render:      render,
		key:         &secure.Key{},
	}
}
