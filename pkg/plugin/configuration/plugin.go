package configuration

import (
	"context"

	"github.com/hashicorp/go-plugin"
	configuration "github.com/wabenet/dodo-core/api/configuration/v1alpha1"
	core "github.com/wabenet/dodo-core/api/core/v1alpha7"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
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
		return nil, dodo.InvalidError{
			Plugin:  p.PluginInfo().GetName(),
			Message: "plugin does not implement Configuration API",
		}
	}

	return &grpcPlugin{Impl: config}, nil
}

type Configuration interface {
	dodo.Plugin

	ListBackdrops() ([]*core.Backdrop, error)
	GetBackdrop(name string) (*core.Backdrop, error)
}

type grpcPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Configuration
}

func (p *grpcPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{configClient: configuration.NewPluginClient(conn)}, nil
}

func (p *grpcPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	configuration.RegisterPluginServer(s, NewGRPCServer(p.Impl))

	return nil
}
