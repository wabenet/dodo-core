package configuration

import (
	api "github.com/dodo-cli/dodo-core/api/v1alpha3"
	dodo "github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const Type pluginType = "configuration"

type pluginType string

func (t pluginType) String() string {
	return string(t)
}

func (t pluginType) GRPCClient() (plugin.Plugin, error) {
	return &grpcPlugin{}, nil
}

func (t pluginType) GRPCServer(p dodo.Plugin) (plugin.Plugin, error) {
	config, ok := p.(Configuration)
	if !ok {
		return nil, dodo.ErrInvalidPlugin{
			Plugin:  p.PluginInfo().Name,
			Message: "plugin does not implement Configuration API",
		}
	}

	return &grpcPlugin{Impl: config}, nil
}

type Configuration interface {
	dodo.Plugin

	ListBackdrops() ([]*api.Backdrop, error)
	GetBackdrop(string) (*api.Backdrop, error)
}

type grpcPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Configuration
}

func (p *grpcPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{configClient: api.NewConfigurationPluginClient(conn)}, nil
}

func (p *grpcPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	api.RegisterConfigurationPluginServer(s, &server{impl: p.Impl})

	return nil
}
