package command

import (
	"github.com/oclaussen/dodo/pkg/plugin/command/proto"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type client struct {
	cmdClient proto.CommandClient
}

func (c *client) GetCommand() (*cobra.Command, error) {
	cmd, err := c.cmdClient.Describe(context.Background(), &proto.CommandQuery{})
	if err != nil {
		return nil, err
	}
	return c.protoToCobra(cmd), nil
}

func (c *client) protoToCobra(in *proto.CommandInfo) *cobra.Command {
	cmd := &cobra.Command{
		Use:              in.Use,
		Short:            in.Short,
		TraverseChildren: in.TraverseChildren,
		SilenceUsage:     in.SilenceUsage,
	}
	if in.ArgsFunc {
		cmd.Args = func(_ *cobra.Command, args []string) error {
			_, err := c.cmdClient.Args(context.Background(), &proto.Arguments{
				Path: in.ExecutePath,
				Args: args,
			})
			return err
		}
	}
	if in.RunFunc {
		cmd.RunE = func(_ *cobra.Command, args []string) error {
			_, err := c.cmdClient.Run(context.Background(), &proto.Arguments{
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
