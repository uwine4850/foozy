## Package mic
Пакет выдповідає за функціональність мікросервісів.

__Connect__
```
func Connect(addr string, onConnect func(conn *grpc.ClientConn, ctx context.Context), opts ...grpc.DialOption) error
```
Підключення до запущеного мікросервіса.