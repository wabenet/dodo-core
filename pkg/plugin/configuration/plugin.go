package configuration

import (
	"github.com/hashicorp/go-plugin"
	dodo "github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/dodo/pkg/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const Type pluginType = "configuration"

type pluginType string

func (t pluginType) String() string {
	return string(t)
}

func (t pluginType) GRPCClient() plugin.Plugin {
	return &grpcPlugin{}
}

func (t pluginType) GRPCServer(p dodo.Plugin) (plugin.Plugin, error) {
	config, ok := p.(Configuration)
	if !ok {
		return nil, dodo.ErrPluginInvalid
	}
	return &grpcPlugin{Impl: config}, nil
}

type Configuration interface {
	Init() error
	UpdateConfiguration(*types.Backdrop) (*types.Backdrop, error)
	Provision(string) error
}

type grpcPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Configuration
}

func (p *grpcPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{configClient: types.NewConfigurationClient(conn)}, nil
}

func (p *grpcPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	types.RegisterConfigurationServer(s, &server{impl: p.Impl})
	return nil
}
