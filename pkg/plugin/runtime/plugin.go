package runtime

import (
	"context"
	"os"

	"github.com/hashicorp/go-plugin"
	core "github.com/wabenet/dodo-core/api/core/v1alpha5"
	runtime "github.com/wabenet/dodo-core/api/runtime/v1alpha1"
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
			Plugin:  p.PluginInfo().Name,
			Message: "plugin does not implement ContainerRuntime API",
		}
	}

	return &grpcPlugin{Impl: rt}, nil
}

type ContainerRuntime interface {
	dodo.Plugin

	ResolveImage(string) (string, error)
	CreateContainer(*core.Backdrop, bool, bool) (string, error)
	StartContainer(string) error
	DeleteContainer(string) error
	ResizeContainer(string, uint32, uint32) error
	KillContainer(string, os.Signal) error
	StreamContainer(string, *dodo.StreamConfig) (*Result, error)
}

type Result struct {
	ExitCode int
}

type grpcPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl ContainerRuntime
}

func (p *grpcPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{runtimeClient: runtime.NewPluginClient(conn)}, nil
}

func (p *grpcPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	runtime.RegisterPluginServer(s, NewGRPCServer(p.Impl))

	return nil
}
