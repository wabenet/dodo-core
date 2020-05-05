package container

import (
	"io"

	"github.com/hashicorp/go-plugin"
	dodoplugin "github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/dodo/pkg/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const PluginType = "containerRuntime"

type ContainerRuntime interface {
	Init() error
	ResolveImage(string) (string, error)
	CreateContainer(*types.Backdrop) (string, error)
	StartContainer(string) error
	RemoveContainer(string) error
	ResizeContainer(string, uint32, uint32) error
	StreamContainer(string, io.Reader, io.Writer) error
}

type Plugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl ContainerRuntime
}

func GetRuntime() (ContainerRuntime, error) {
	for _, p := range dodoplugin.GetPlugins(PluginType) {
		log.Debug("trying plugin")
		if rt, ok := p.(ContainerRuntime); ok {
			log.Debug("ok")
			return rt, nil
		}
	}
	return nil, errors.New("Could not find any container runtime")

}

func init() {
	dodoplugin.RegisterPluginClient(PluginType, &Plugin{})
}

func (p *Plugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{runtimeClient: types.NewContainerRuntimeClient(conn)}, nil
}

func (p *Plugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	types.RegisterContainerRuntimeServer(s, &server{impl: p.Impl})
	return nil
}
