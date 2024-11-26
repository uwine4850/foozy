package manager

import (
	"sync"

	"github.com/uwine4850/foozy/pkg/interfaces"
)

type OneTimeData struct {
	userContext sync.Map
	slugParams  map[string]string
}

func NewManagerData() *OneTimeData {
	return &OneTimeData{}
}

func (m *OneTimeData) New() (interface{}, error) {
	return &OneTimeData{}, nil
}

// SetSlugParams sets the slug parameters.
func (m *OneTimeData) SetSlugParams(params map[string]string) {
	m.slugParams = params
}

// GetSlugParams returns the parameter by key. If the key is not found returns false.
func (m *OneTimeData) GetSlugParams(key string) (string, bool) {
	res, ok := m.slugParams[key]
	return res, ok
}

// SetUserContext sets the user context.
// This context is used only as a means of passing information between handlers.
func (m *OneTimeData) SetUserContext(key string, value interface{}) {
	m.userContext.Store(key, value)
}

// GetUserContext getting the user context.
func (m *OneTimeData) GetUserContext(key string) (any, bool) {
	value, ok := m.userContext.Load(key)
	return value, ok
}

// DelUserContext deletes a user context by key.
func (m *OneTimeData) DelUserContext(key string) {
	m.userContext.Delete(key)
}

// CreateAndSetNewManagerData сreates and sets a new OneTimeData instance into the manager.
func CreateAndSetNewManagerData(manager interfaces.IManager) error {
	otd := manager.OneTimeData()
	newOtd, err := otd.New()
	if err != nil {
		return err
	}
	manager.SetOneTimeData(newOtd.(interfaces.IManagerOneTimeData))
	return nil
}
