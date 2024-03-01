package mic

import (
	"context"
	"google.golang.org/grpc"
)

func Connect(addr string, onConnect func(conn *grpc.ClientConn, ctx context.Context), opts ...grpc.DialOption) error {
	o := opts
	if opts == nil {
		o = []grpc.DialOption{grpc.WithInsecure()}
	} else {
		o = opts
	}
	conn, err := grpc.Dial(addr, o...)
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx := context.Background()
	onConnect(conn, ctx)
	return nil
}
