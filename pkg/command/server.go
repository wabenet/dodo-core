package command

import (
	"strings"

	"github.com/oclaussen/dodo/pkg/types"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type server struct {
	impl Command
}

func (s *server) Describe(ctx context.Context, _ *types.Empty) (*types.Commands, error) {
	cmds, err := s.impl.GetCommands()
	if err != nil {
		return nil, err
	}
	result := map[string]*types.CommandInfo{}
	for name, cmd := range cmds {
		result[name] = s.cobraToProto(cmd)
	}
	return &types.Commands{Commands: result}, nil
}

func (s *server) Run(ctx context.Context, args *types.CommandArguments) (*types.Empty, error) {
	cmds, err := s.impl.GetCommands()
	if err != nil {
		return nil, err
	}
        cmd := cmds[args.Path[0]]
	subCmd, _, err := cmd.Find(args.Path[1:])
	if err != nil {
		return nil, err
	}
	if err = subCmd.RunE(subCmd, args.Args); err != nil {
		return nil, err
	}
	return &types.Empty{}, nil
}

func (s *server) Args(ctx context.Context, args *types.CommandArguments) (*types.Empty, error) {
	cmds, err := s.impl.GetCommands()
	if err != nil {
		return nil, err
	}
        cmd := cmds[args.Path[0]]
	subCmd, _, err := cmd.Find(args.Path[1:])
	if err != nil {
		return nil, err
	}
	if err = subCmd.Args(subCmd, args.Args); err != nil {
		return nil, err
	}
	return &types.Empty{}, nil
}

func (s *server) cobraToProto(in *cobra.Command) *types.CommandInfo {
	subcommands := []*types.CommandInfo{}
	for _, sub := range in.Commands() {
		subcommands = append(subcommands, s.cobraToProto(sub))
	}
	cmd := &types.CommandInfo{
		Use:              in.Use,
		Short:            in.Short,
		TraverseChildren: in.TraverseChildren,
		SilenceUsage:     in.SilenceUsage,
		Subcommands:      subcommands,
	}
	cmd.ExecutePath = strings.Split(in.CommandPath(), " ")
	cmd.ArgsFunc = (in.Args != nil)
	cmd.RunFunc = (in.RunE != nil)
	return cmd
}
