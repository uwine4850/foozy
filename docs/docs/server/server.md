## package server
The package is responsible for the operation of the framework server.

### Server
An object that stores server data and launches the server.

#### NewServer
Creating a new server instance.<br>
The `rcors` argument is responsible for launching `CORS` together with the `server`. The argument can be nil.<br>
`CORS` is implemented using [https://github.com/rs/cors](https://github.com/rs/cors).
```golang
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
```

#### Server.Start
Starts listening for connections on the selected address.
```golang
func (s *Server) Start() error {
	fmt.Printf("Server start on %s", s.addr)
	err := s.serv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
```

#### Server.GetServ
Returns an `*http.Server` object.
```golang
func (s *Server) GetServ() *http.Server {
	return s.serv
}
```

#### Server.Stop
Stopping server listening.
```golang
func (s *Server) Stop() error {
	if err := s.serv.Close(); err != nil {
		return err
	}
	println("Server stopped.")
	return nil
}
```

#### MicServer
Server implementation object for microservices.

#### MicServer.Start
Starts listening for connections on the selected address.
```golang
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
```

#### FoozyAndMic
Launching two servers simultaneously:

* Server
* MicServer

```golang
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
```

#### WaitStartServer
Waits for the server to start at the selected address for the specified amount of time.<br>
If the connection did not appear during this time, it returns an error.<br>
__IMPORTANT:__ Either the server or this function must be run in a goroutine, otherwise it may not work.
```golang
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
```