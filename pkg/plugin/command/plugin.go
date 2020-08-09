package command

import (
	"github.com/hashicorp/go-plugin"
	dodo "github.com/oclaussen/dodo/pkg/plugin"
	"github.com/spf13/cobra"
)

// TODO: Implement cobra commands via GRPC

const Type pluginType = "command"

type pluginType string

func (t pluginType) String() string {
	return string(t)
}

func (t pluginType) GRPCClient() (plugin.Plugin, error) {
	return nil, dodo.ErrPluginNotImplemented
}

func (t pluginType) GRPCServer(p dodo.Plugin) (plugin.Plugin, error) {
	return nil, dodo.ErrPluginNotImplemented
}

type Command interface {
	Init() error
	Name() string
	GetCobraCommand() *cobra.Command
}
