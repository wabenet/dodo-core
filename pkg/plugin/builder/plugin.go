package builder

import (
	"context"

	"github.com/hashicorp/go-plugin"
	build "github.com/wabenet/dodo-core/api/build/v1alpha1"
	core "github.com/wabenet/dodo-core/api/core/v1alpha6"
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
			Plugin:  p.PluginInfo().GetName(),
			Message: "plugin does not implement ImageBuilder API",
		}
	}

	return &grpcPlugin{Impl: rt}, nil
}

type ImageBuilder interface {
	dodo.Plugin

	CreateImage(info *core.BuildInfo, streamConfig *dodo.StreamConfig) (string, error)
}

type grpcPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl ImageBuilder
}

func (p *grpcPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{builderClient: build.NewPluginClient(conn)}, nil
}

func (p *grpcPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	build.RegisterPluginServer(s, NewGRPCServer(p.Impl))

	return nil
}
