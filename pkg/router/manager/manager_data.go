package manager

import (
	"sync"
)

type ManagerData struct {
	userContext sync.Map
	slugParams  map[string]string
}

func NewManagerData() *ManagerData {
	return &ManagerData{}
}

// SetSlugParams sets the slug parameters.
func (m *ManagerData) SetSlugParams(params map[string]string) {
	m.slugParams = params
}

// GetSlugParams returns the parameter by key. If the key is not found returns false.
func (m *ManagerData) GetSlugParams(key string) (string, bool) {
	res, ok := m.slugParams[key]
	return res, ok
}

// SetUserContext sets the user context.
// This context is used only as a means of passing information between handlers.
func (m *ManagerData) SetUserContext(key string, value interface{}) {
	m.userContext.Store(key, value)
}

// GetUserContext getting the user context.
func (m *ManagerData) GetUserContext(key string) (any, bool) {
	value, ok := m.userContext.Load(key)
	return value, ok
}

// DelUserContext deletes a user context by key.
func (m *ManagerData) DelUserContext(key string) {
	m.userContext.Delete(key)
}
