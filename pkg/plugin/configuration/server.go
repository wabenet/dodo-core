package configuration

import (
	"github.com/oclaussen/dodo/pkg/plugin/configuration/proto"
	"github.com/oclaussen/dodo/pkg/types"
	"golang.org/x/net/context"
)

type server struct {
	impl Configuration
}

func (s *server) GetClientOptions(_ context.Context, request *proto.Backdrop) (*proto.ClientOptions, error) {
	opts, err := s.impl.GetClientOptions(request.Name)
	if err != nil {
		return nil, err
	}
	return &proto.ClientOptions{
		Version:  opts.Version,
		Host:     opts.Host,
		CaFile:   opts.CAFile,
		CertFile: opts.CertFile,
		KeyFile:  opts.KeyFile,
	}, nil
}

func (s *server) UpdateConfiguration(_ context.Context, backdrop *types.Backdrop) (*types.Backdrop, error) {
	return s.impl.UpdateConfiguration(backdrop)
}

func (s *server) Provision(_ context.Context, request *proto.Container) (*proto.Empty, error) {
	return &proto.Empty{}, s.impl.Provision(request.Id)
}
