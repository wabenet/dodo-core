package configuration

import (
	"github.com/hashicorp/go-plugin"
	dodoplugin "github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/dodo/pkg/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const PluginType = "configuration"

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

func init() {
	dodoplugin.RegisterPluginClient(PluginType, &Plugin{})
}

func (p *Plugin) GRPCClient(_ context.Context, _ *plugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &client{configClient: types.NewConfigurationClient(conn)}, nil
}

func (p *Plugin) GRPCServer(_ *plugin.GRPCBroker, s *grpc.Server) error {
	types.RegisterConfigurationServer(s, &server{impl: p.Impl})
	return nil
}
