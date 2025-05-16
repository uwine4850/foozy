package manager

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/secure"
)

// Manager structure for managing router handlers.
// The main function of [Manager], is to store and transfer different data.
// IMPORTANT: a separate instance of [Manager] must be created for each new http request.
// This is very important for data security, as the object stores data directly in itself.
type Manager struct {
	oneTimeData  interfaces.IManagerOneTimeData
	render       interfaces.IRender
	key          interfaces.IKey
	databasePool interfaces.IDatabasePool
}

// New creates a new instance of [Manager] with some of the old settings.
// This is a very important method because it creates a new instance of [Manager],
// but retains the static data. If a new instance is created without this method,
// the router may not work properly or may cause serious problems.
func (m *Manager) New() (interface{}, error) {
	newOTD, err := m.oneTimeData.New()
	if err != nil {
		return nil, err
	}
	var newRender interfaces.IRender
	if m.render != nil {
		_newRender, err := m.render.New()
		if err != nil {
			return nil, err
		}
		newRender = _newRender.(interfaces.IRender)
	}

	return &Manager{
		oneTimeData:  newOTD.(interfaces.IManagerOneTimeData),
		render:       newRender,
		key:          m.key,
		databasePool: m.databasePool,
	}, nil
}

func (m *Manager) Render() interfaces.IRender {
	return m.render
}

func (m *Manager) OneTimeData() interfaces.IManagerOneTimeData {
	return m.oneTimeData
}

func (m *Manager) Key() interfaces.IKey {
	return m.key
}

func (m *Manager) Database() interfaces.IDatabasePool {
	return m.databasePool
}

func NewManager(otd interfaces.IManagerOneTimeData, render interfaces.IRender, databasePool interfaces.IDatabasePool) *Manager {
	return &Manager{
		oneTimeData:  otd,
		render:       render,
		key:          secure.NewKey(),
		databasePool: databasePool,
	}
}
