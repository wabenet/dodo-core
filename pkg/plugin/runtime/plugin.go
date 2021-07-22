package runtime

import (
	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	dodo "github.com/dodo-cli/dodo-core/pkg/plugin"
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
	dodo.Plugin

	ResolveImage(string) (string, error)
	CreateContainer(*api.Backdrop, bool, bool) (string, error)
	StartContainer(string) error
	DeleteContainer(string) error
	ResizeContainer(string, uint32, uint32) error
	StreamContainer(string, *dodo.StreamConfig) error
}

type Result struct {
	ExitCode int64
	Message  string
}

func (e Result) Error() string {
	return e.Message
}

type grpcPlugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl ContainerRuntime
}

func (p *grpcPlugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{runtimeClient: api.NewRuntimePluginClient(conn)}, nil
}

func (p *grpcPlugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	api.RegisterRuntimePluginServer(s, &server{impl: p.Impl})
	return nil
}
