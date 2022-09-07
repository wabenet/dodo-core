package test

import (
	"context"
	"fmt"
	"log"
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
		if err := s.Serve(lis); err != nil {
			log.Println(err.Error())
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			conn, err := lis.Dial()
			if err != nil {
				return nil, fmt.Errorf("could not dial bufconn: %w", err)
			}

			return conn, nil
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, func() {}, fmt.Errorf("could not connect to bufconn: %w", err)
	}

	c, err := grpcClient(t, conn)
	if err != nil {
		conn.Close()

		return nil, func() {}, err
	}

	if _, err := c.Init(); err != nil {
		conn.Close()

		return nil, func() {}, fmt.Errorf("could not init client: %w", err)
	}

	return c, func() {
		conn.Close()
		c.Cleanup()
	}, nil
}

func grpcServer(t dodo.Type, p dodo.Plugin) (*grpc.Server, error) {
	pl, err := t.GRPCServer(p)
	if err != nil {
		return nil, fmt.Errorf("could not setup grpc server: %w", err)
	}

	gp, ok := pl.(plugin.GRPCPlugin)
	if !ok {
		return nil, dodo.InvalidError{}
	}

	s := grpc.NewServer()

	if err := gp.GRPCServer(nil, s); err != nil {
		return nil, fmt.Errorf("could not start grpc server: %w", err)
	}

	return s, nil
}

func grpcClient(t dodo.Type, conn *grpc.ClientConn) (dodo.Plugin, error) {
	p, err := t.GRPCClient()
	if err != nil {
		return nil, fmt.Errorf("could not setup grpc client: %w", err)
	}

	gp, ok := p.(plugin.GRPCPlugin)
	if !ok {
		return nil, dodo.InvalidError{}
	}

	c, err := gp.GRPCClient(context.Background(), nil, conn)
	if err != nil {
		return nil, fmt.Errorf("could not start grpc client: %w", err)
	}

	b, ok := c.(dodo.Plugin)
	if !ok {
		return nil, dodo.InvalidError{}
	}

	return b, nil
}
