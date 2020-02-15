package command

import (
	"strings"

	"github.com/oclaussen/dodo/pkg/plugin/command/proto"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

type server struct {
	impl Command
}

func (s *server) Describe(ctx context.Context, _ *proto.CommandQuery) (*proto.CommandInfo, error) {
	cmd, err := s.impl.GetCommand()
	if err != nil {
		return nil, err
	}
	return s.cobraToProto(cmd), nil
}

func (s *server) Run(ctx context.Context, args *proto.Arguments) (*proto.Result, error) {
	cmd, err := s.impl.GetCommand()
	if err != nil {
		return nil, err
	}
	subCmd, _, err := cmd.Find(args.Path[1:])
	if err != nil {
		return nil, err
	}
	if err = subCmd.RunE(subCmd, args.Args); err != nil {
		return nil, err
	}
	return &proto.Result{}, nil
}

func (s *server) Args(ctx context.Context, args *proto.Arguments) (*proto.Result, error) {
	cmd, err := s.impl.GetCommand()
	if err != nil {
		return nil, err
	}
	subCmd, _, err := cmd.Find(args.Path[1:])
	if err != nil {
		return nil, err
	}
	if err = subCmd.Args(subCmd, args.Args); err != nil {
		return nil, err
	}
	return &proto.Result{}, nil
}

func (s *server) cobraToProto(in *cobra.Command) *proto.CommandInfo {
	subcommands := []*proto.CommandInfo{}
	for _, sub := range in.Commands() {
		subcommands = append(subcommands, s.cobraToProto(sub))
	}
	cmd := &proto.CommandInfo{
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
