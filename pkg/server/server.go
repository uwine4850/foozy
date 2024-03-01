package server

import (
	"fmt"
	"github.com/uwine4850/foozy/pkg/router"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"sync"
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

type MicServer struct {
	Network    string
	Address    string
	GrpcServer *grpc.Server
}

func NewMicServer(network string, address string, grpcServer *grpc.Server) *MicServer {
	return &MicServer{Network: network, Address: address, GrpcServer: grpcServer}
}

func (mc *MicServer) Start() error {
	lis, err := net.Listen(mc.Network, mc.Address)
	if err != nil {
		return err
	}
	log.Printf("Mic server start on %v", lis.Addr())
	if err := mc.GrpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func FoozyAndMic(fserver *Server, micServer *MicServer, onError func(err error)) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := fserver.Start()
		if err != nil {
			onError(err)
		}
	}()
	wg.Add(1)
	go func() {
		err := micServer.Start()
		if err != nil {
			onError(err)
		}
	}()
	wg.Wait()
}
