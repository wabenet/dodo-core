package runtime

import (
	"context"

	"github.com/hashicorp/go-plugin"
	pluginapi "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/runtime/v1alpha2"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
	"google.golang.org/grpc"
)

const Type pluginType = "runtime"

type pluginType string

func (t pluginType) String() string {
	return string(t)
}

func (t pluginType) GRPCClient() (plugin.Plugin, error) {
	return &grpcPlugin{}, nil
}

func (t pluginType) GRPCServer(p dodo.Plugin) (plugin.Plugin, error) {
	rt, ok := p.(ContainerRuntime)
	if !ok {
		return nil, dodo.InvalidError{
			PluginID: p.Metadata().ID,
			Message:  "plugin does not implement ContainerRuntime API",
		}
	}

	return &grpcPlugin{Impl: rt}, nil
}

type grpcPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl ContainerRuntime
}

func (p *grpcPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return NewGRPCClient(conn), nil
}

func (p *grpcPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	impl := NewGRPCServer(p.Impl)

	pluginapi.RegisterPluginServer(s, impl)
	pluginapi.RegisterOutputStreamingPluginServer(s, impl)
	pluginapi.RegisterInputStreamingPluginServer(s, impl)
	api.RegisterRuntimePluginServer(s, NewGRPCServer(p.Impl))

	return nil
}
