package command

import (
	"github.com/hashicorp/go-plugin"
	"github.com/oclaussen/dodo/pkg/plugin/command/proto"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Command interface {
	GetCommand() (*cobra.Command, error)
}

type Plugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Command
}

func (p *Plugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{cmdClient: proto.NewCommandClient(conn)}, nil
}

func (p *Plugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterCommandServer(s, &server{impl: p.Impl})
	return nil
}
