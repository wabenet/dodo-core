package command

import (
	"github.com/oclaussen/dodo/pkg/types"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type client struct {
	cmdClient types.CommandClient
}

func (c *client) GetCommands() (map[string]*cobra.Command, error) {
	response, err := c.cmdClient.Describe(context.Background(), &types.Empty{})
	if err != nil {
		return nil, err
	}
	cmds := map[string]*cobra.Command{}
	for name, cmd := range response.Commands {
		cmds[name] = c.protoToCobra(cmd)
	}
	return cmds, nil
}

func (c *client) protoToCobra(in *types.CommandInfo) *cobra.Command {
	cmd := &cobra.Command{
		Use:              in.Use,
		Short:            in.Short,
		TraverseChildren: in.TraverseChildren,
		SilenceUsage:     in.SilenceUsage,
	}
	if in.ArgsFunc {
		cmd.Args = func(_ *cobra.Command, args []string) error {
			_, err := c.cmdClient.Args(context.Background(), &types.CommandArguments{
				Path: in.ExecutePath,
				Args: args,
			})
			return err
		}
	}
	if in.RunFunc {
		cmd.RunE = func(_ *cobra.Command, args []string) error {
			_, err := c.cmdClient.Run(context.Background(), &types.CommandArguments{
				Path: in.ExecutePath,
				Args: args,
			})
			return err
		}
	}
	for _, sub := range in.Subcommands {
		cmd.AddCommand(c.protoToCobra(sub))
	}
	return cmd
}
