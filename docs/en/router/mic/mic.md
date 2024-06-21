## Package mic
The package is responsible for the functionality of microservices.

__Connect__
```
func Connect(addr string, onConnect func(conn *grpc.ClientConn, ctx context.Context), opts ...grpc.DialOption) error
```
Connecting to a running microservice.