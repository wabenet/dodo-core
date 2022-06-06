package test

import (
	"context"
	"net"

	"github.com/hashicorp/go-plugin"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func GRPCWrapPlugin(t dodo.Type, p dodo.Plugin) (dodo.Plugin, func(), error) {
	lis := bufconn.Listen(bufSize)

	s, err := grpcServer(t, p)
	if err != nil {
		return nil, func() {}, err
	}

	go func() {
		s.Serve(lis)
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, func() {}, err
	}

	cleanup := func() { conn.Close() }

	c, err := grpcClient(t, conn)
	if err != nil {
		return nil, cleanup, err
	}

	return c, cleanup, nil
}

func grpcServer(t dodo.Type, p dodo.Plugin) (*grpc.Server, error) {
	pl, err := t.GRPCServer(p)
	if err != nil {
		return nil, err
	}

	gp, ok := pl.(plugin.GRPCPlugin)
	if !ok {
		return nil, dodo.InvalidError{}
	}

	s := grpc.NewServer()

	if err := gp.GRPCServer(nil, s); err != nil {
		return nil, err
	}

	return s, nil
}

func grpcClient(t dodo.Type, conn *grpc.ClientConn) (dodo.Plugin, error) {
	p, err := t.GRPCClient()
	if err != nil {
		return nil, err
	}

	gp, ok := p.(plugin.GRPCPlugin)
	if !ok {
		return nil, dodo.InvalidError{}
	}

	c, err := gp.GRPCClient(context.Background(), nil, conn)
	if err != nil {
		return nil, err
	}

	b, ok := c.(dodo.Plugin)
	if !ok {
		return nil, dodo.InvalidError{}
	}

	return b, nil
}
