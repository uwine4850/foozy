package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/rs/cors"
	"github.com/uwine4850/foozy/pkg/router"
	"google.golang.org/grpc"
)

type Server struct {
	router *router.Router
	addr   string
	serv   *http.Server
}

// rcors - CORS data, which is set using the cors.Options library at github.com/rs/cors.
// These settings will be applied to all requests.
func NewServer(addr string, router *router.Router, rcors *cors.Options) *Server {
	s := &Server{router: router, addr: addr}
	var handler http.Handler = router

	if rcors != nil {
		handler = cors.New(*rcors).Handler(handler)
	}

	s.serv = &http.Server{
		Addr:    s.addr,
		Handler: handler,
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
	if err := s.serv.Close(); err != nil {
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

type ErrResponseTimedOut struct {
	Address string
}

func (e ErrResponseTimedOut) Error() string {
	return fmt.Sprintf("the response timed out, address: %s", e.Address)
}

// WaitStartServer waits for the server to start at the selected address for the specified amount of time.
// If the connection did not appear during this time, it returns an error.
// IMPORTANT: Either the server or this function must be run in a goroutine, otherwise it may not work.
func WaitStartServer(addr string, waitTimeSec int) error {
	var outErr error
	for i := 0; i < waitTimeSec; i++ {
		time.Sleep(1 * time.Second)
		conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
		if conn != nil {
			conn.Close()
		}
		if err != nil {
			outErr = ErrResponseTimedOut{addr}
		} else {
			outErr = nil
			break
		}
	}
	return outErr
}
