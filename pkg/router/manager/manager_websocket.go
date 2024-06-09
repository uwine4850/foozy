package manager

import "github.com/uwine4850/foozy/pkg/interfaces"

type ManagerWebsocket struct {
	websocket interfaces.IWebsocket
}

func NewManagerWebsocket() *ManagerWebsocket {
	return &ManagerWebsocket{}
}

// CurrentWebsocket get an instance for the websocket connection.
// Works only in the "Ws" handler.
func (m *ManagerWebsocket) CurrentWebsocket() interfaces.IWebsocket {
	return m.websocket
}

// SetWebsocket sets the websocket interface.
func (m *ManagerWebsocket) SetWebsocket(websocket interfaces.IWebsocket) {
	m.websocket = websocket
}
