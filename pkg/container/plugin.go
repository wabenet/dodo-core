package container

import (
	"fmt"
	"io"

	"github.com/hashicorp/go-plugin"
	dodo "github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/dodo/pkg/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const PluginType = "containerRuntime"

type ContainerRuntime interface {
	Init() error
	ResolveImage(string) (string, error)
	CreateContainer(*types.Backdrop, bool, bool) (string, error)
	StartContainer(string) error
	RemoveContainer(string) error
	ResizeContainer(string, uint32, uint32) error
	StreamContainer(string, io.Reader, io.Writer) error
}

type Plugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl ContainerRuntime
}

func RegisterPlugin() {
	dodo.RegisterPluginClient(PluginType, &Plugin{})
}

func GetRuntime() (ContainerRuntime, error) {
	for _, p := range dodo.GetPlugins(PluginType) {
		if rt, ok := p.(ContainerRuntime); ok {
			return rt, nil
		}
	}

	return nil, fmt.Errorf("could not find container runtime: %w", dodo.ErrNoValidPluginFound)
}

func (p *Plugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{runtimeClient: types.NewContainerRuntimeClient(conn)}, nil
}

func (p *Plugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	types.RegisterContainerRuntimeServer(s, &server{impl: p.Impl})
	return nil
}
