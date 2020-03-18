package command

import (
	"github.com/hashicorp/go-plugin"
	dodoplugin "github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/dodo/pkg/types"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const PluginType = "command"

type Command interface {
	GetCommands() (map[string]*cobra.Command, error)
}

type Plugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Command
}

func init() {
	dodoplugin.RegisterPluginClient(PluginType, &Plugin{})
}

func (p *Plugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{cmdClient: types.NewCommandClient(conn)}, nil
}

func (p *Plugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	types.RegisterCommandServer(s, &server{impl: p.Impl})
	return nil
}
