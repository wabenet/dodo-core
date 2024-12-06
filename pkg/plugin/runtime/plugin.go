package runtime

import (
	"context"
	"os"

	"github.com/hashicorp/go-plugin"
	core "github.com/wabenet/dodo-core/api/core/v1alpha6"
	runtime "github.com/wabenet/dodo-core/api/runtime/v1alpha2"
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
			Plugin:  p.PluginInfo().GetName(),
			Message: "plugin does not implement ContainerRuntime API",
		}
	}

	return &grpcPlugin{Impl: rt}, nil
}

type ContainerRuntime interface {
	dodo.Plugin

	ResolveImage(spec string) (string, error)
	CreateContainer(backdrop *core.Backdrop, tty, stdio bool) (string, error)
	StartContainer(id string) error
	DeleteContainer(id string) error
	ResizeContainer(id string, height, width uint32) error
	KillContainer(id string, signal os.Signal) error
	StreamContainer(id string, streamConfig *dodo.StreamConfig) (*Result, error)
	CreateVolume(name string) error
	DeleteVolume(name string) error
	WriteFile(name, path string, contents []byte) error
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
