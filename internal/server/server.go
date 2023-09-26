package server

import (
	"foozy/internal/router"
	"net/http"
)

type Server struct {
	router router.IRouter
	addr   string
}

func NewServer(addr string, router router.IRouter) *Server {
	return &Server{router: router, addr: addr}
}

func (s *Server) Start() error {
	err := http.ListenAndServe(s.addr, s.router.GetMux())
	if err != nil {
		return err
	}
	return nil
}
