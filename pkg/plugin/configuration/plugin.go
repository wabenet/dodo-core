package configuration

import (
	"github.com/hashicorp/go-plugin"
	"github.com/oclaussen/dodo/pkg/plugin/configuration/proto"
	"github.com/oclaussen/dodo/pkg/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Configuration interface {
	GetClientOptions(string) (*ClientOptions, error)
	UpdateConfiguration(*types.Backdrop) (*types.Backdrop, error)
	Provision(string) error
}

type ClientOptions struct {
	Version  string
	Host     string
	CAFile   string
	CertFile string
	KeyFile  string
}

type Plugin struct {
	plugin.NetRPCUnsupportedPlugin
	Impl Configuration
}

func (p *Plugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{configClient: proto.NewConfigurationClient(conn)}, nil
}

func (p *Plugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterConfigurationServer(s, &server{impl: p.Impl})
	return nil
}
