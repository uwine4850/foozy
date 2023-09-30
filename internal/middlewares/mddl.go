package middlewares

import (
	"github.com/uwine4850/foozy/internal/interfaces"
	"github.com/uwine4850/foozy/internal/utils"
	"net/http"
	"sync"
)

type MddlFunc func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)

type Middleware struct {
	preHandlerMiddlewares map[int]MddlFunc
	preHandlerId          []int
	Context               sync.Map
	mError                error
}

func NewMiddleware() *Middleware {
	return &Middleware{preHandlerMiddlewares: make(map[int]MddlFunc)}
}

func (m *Middleware) PreHandlerMddl(id int, fn func(w http.ResponseWriter, r *http.Request, manager interfaces.IManager)) {
	if !utils.SliceContains(m.preHandlerId, id) {
		m.preHandlerId = append(m.preHandlerId, id)
		m.preHandlerMiddlewares[id] = fn
	} else {
		m.mError = &ErrIdAlreadyExist{id}
	}
}

func (m *Middleware) RunPreMddl(w http.ResponseWriter, r *http.Request, manager interfaces.IManager) error {
	if m.mError != nil {
		return m.mError
	}
	for _, handlerFunc := range m.preHandlerMiddlewares {
		handlerFunc(w, r, manager)
	}
	return nil
}
