package run

import (
	"github.com/dodo/dodo-core/pkg/plugin"
	"github.com/dodo/dodo-core/pkg/plugin/command"
	"github.com/spf13/cobra"
)

const name = "run"

type Command struct {
	cmd *cobra.Command
}

func (p *Command) Type() plugin.Type {
	return command.Type
}

func (p *Command) Init() error {
	p.cmd = NewCommand()
	return nil
}

func (p *Command) Name() string {
	return name
}

func (p *Command) GetCobraCommand() *cobra.Command {
	return p.cmd
}
