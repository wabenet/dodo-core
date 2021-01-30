package command

import (
	dodo "github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/hashicorp/go-plugin"
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
	dodo.Plugin

	GetCobraCommand() *cobra.Command
}
