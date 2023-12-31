package server

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/router"
	"log"
	"net/http"
)

type Server struct {
	router *router.Router
	addr   string
	serv   *http.Server
}

func NewServer(addr string, router *router.Router) *Server {
	s := &Server{router: router, addr: addr}
	s.serv = &http.Server{
		Addr:    s.addr,
		Handler: s.router.GetMux(),
	}
	return s
}

// Start starts the server.
func (s *Server) Start() error {
	log.Println(fmt.Sprintf("Server start on %s", s.addr))
	err := s.serv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) GetServ() *http.Server {
	return s.serv
}

// Stop stops the server.
func (s *Server) Stop() error {
	err := s.serv.Shutdown(nil)
	if err != nil {
		return err
	}
	println("Server stopped.")
	return nil
}
