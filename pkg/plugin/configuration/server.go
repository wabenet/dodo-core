package configuration

import (
	"fmt"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

type server struct {
	impl Configuration
}

func (s *server) GetPluginInfo(_ context.Context, _ *empty.Empty) (*api.PluginInfo, error) {
	return s.impl.PluginInfo()
}

func (s *server) ListBackdrops(_ context.Context, _ *empty.Empty) (*api.ListBackdropsResponse, error) {
	backdrops, err := s.impl.ListBackdrops()
	if err != nil {
		return nil, fmt.Errorf("could not list backdrops: %w", err)
	}

	return &api.ListBackdropsResponse{Backdrops: backdrops}, nil
}

func (s *server) GetBackdrop(_ context.Context, request *api.GetBackdropRequest) (*api.Backdrop, error) {
	return s.impl.GetBackdrop(request.Alias)
}
