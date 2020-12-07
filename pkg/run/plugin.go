package run

import (
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/command"
	"github.com/dodo-cli/dodo-core/pkg/types"
	"github.com/spf13/cobra"
)

const name = "run"

type Command struct {
	cmd *cobra.Command
}

func (p *Command) Type() plugin.Type {
	return command.Type
}

func (p *Command) Init() (*types.PluginInfo, error) {
	p.cmd = NewCommand()
	return &types.PluginInfo{Name: name}, nil
}

func (p *Command) GetCobraCommand() *cobra.Command {
	return p.cmd
}
