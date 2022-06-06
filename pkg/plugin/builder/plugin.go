package builder

import (
	"context"

	"github.com/hashicorp/go-plugin"
	api "github.com/wabenet/dodo-core/api/v1alpha3"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
	"google.golang.org/grpc"
)

const Type pluginType = "builder"

type pluginType string

func (t pluginType) String() string {
	return string(t)
}

func (t pluginType) GRPCClient() (plugin.Plugin, error) {
	return &grpcPlugin{}, nil
}

func (t pluginType) GRPCServer(p dodo.Plugin) (plugin.Plugin, error) {
	rt, ok := p.(ImageBuilder)
	if !ok {
		return nil, dodo.InvalidError{
			Plugin:  p.PluginInfo().Name,
			Message: "plugin does not implement ImageBuilder API",
		}
	}

	return &grpcPlugin{Impl: rt}, nil
}

type ImageBuilder interface {
	dodo.Plugin

	CreateImage(*api.BuildInfo, *dodo.StreamConfig) (string, error)
}

type grpcPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl ImageBuilder
}

func (p *grpcPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{builderClient: api.NewBuilderPluginClient(conn)}, nil
}

func (p *grpcPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	api.RegisterBuilderPluginServer(s, NewGRPCServer(p.Impl))

	return nil
}
