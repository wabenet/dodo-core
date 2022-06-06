package command

import (
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
)

// TODO: Implement cobra commands via GRPC

const Type pluginType = "command"

type pluginType string

func (t pluginType) String() string {
	return string(t)
}

func (t pluginType) GRPCClient() (plugin.Plugin, error) {
	return nil, dodo.InvalidError{
		Message: "command plugin is not implemented via grpc",
	}
}

func (t pluginType) GRPCServer(p dodo.Plugin) (plugin.Plugin, error) {
	return nil, dodo.InvalidError{
		Message: "command plugin is not implemented via grpc",
	}
}

type Command interface {
	dodo.Plugin

	GetCobraCommand() *cobra.Command
}
