package configuration

import (
	"github.com/dodo/dodo-core/pkg/types"
	"golang.org/x/net/context"
)

type server struct {
	impl Configuration
}

func (s *server) Init(_ context.Context, _ *types.Empty) (*types.Empty, error) {
	return &types.Empty{}, s.impl.Init()
}

func (s *server) UpdateConfiguration(_ context.Context, backdrop *types.Backdrop) (*types.Backdrop, error) {
	return s.impl.UpdateConfiguration(backdrop)
}

func (s *server) Provision(_ context.Context, request *types.ContainerId) (*types.Empty, error) {
	return &types.Empty{}, s.impl.Provision(request.Id)
}
