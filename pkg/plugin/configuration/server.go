package configuration

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/wabenet/dodo-core/api/v1alpha3"
)

type server struct {
	impl Configuration
}

func NewGRPCServer(impl Configuration) api.ConfigurationPluginServer {
	return &server{impl: impl}
}

func (s *server) GetPluginInfo(_ context.Context, _ *empty.Empty) (*api.PluginInfo, error) {
	return s.impl.PluginInfo(), nil
}

func (s *server) InitPlugin(_ context.Context, _ *empty.Empty) (*api.InitPluginResponse, error) {
	config, err := s.impl.Init()
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return &api.InitPluginResponse{Config: config}, nil
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
