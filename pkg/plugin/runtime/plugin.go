package runtime

import (
	"io"

	dodo "github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/types"
	"github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
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
		return nil, dodo.ErrPluginInvalid
	}

	return &grpcPlugin{Impl: rt}, nil
}

type ContainerRuntime interface {
	Init() (*types.PluginInfo, error)
	ResolveImage(string) (string, error)
	CreateContainer(*types.Backdrop, bool, bool) (string, error)
	StartContainer(string) error
	RemoveContainer(string) error
	ResizeContainer(string, uint32, uint32) error
	StreamContainer(string, io.Reader, io.Writer, uint32, uint32) error
}

type grpcPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl ContainerRuntime
}

func (p *grpcPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{runtimeClient: types.NewContainerRuntimeClient(conn)}, nil
}

func (p *grpcPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	types.RegisterContainerRuntimeServer(s, &server{impl: p.Impl})
	return nil
}
