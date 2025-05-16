package manager

import (
	"fmt"
	"sync"
)

// OneTimeData one-time manager data.
// This object stores temporary data for the router handler.
// The data goes through a full request cycle and can be used between mddlewares, for example.
//
// You should not store permanent data here, only temporary data for each request.
type OneTimeData struct {
	userContext sync.Map
	slugParams  map[string]string
}

func NewOneTimeData() *OneTimeData {
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
	m.userContext.Range(func(key, value any) bool {
		fmt.Println(key)
		return true
	})
	value, ok := m.userContext.Load(key)
	return value, ok
}

// DelUserContext deletes a user context by key.
func (m *OneTimeData) DelUserContext(key string) {
	m.userContext.Delete(key)
}
