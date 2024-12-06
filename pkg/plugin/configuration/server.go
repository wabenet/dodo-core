package configuration

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	configuration "github.com/wabenet/dodo-core/api/configuration/v1alpha1"
	core "github.com/wabenet/dodo-core/api/core/v1alpha6"
)

type server struct {
	impl Configuration
}

func NewGRPCServer(impl Configuration) configuration.PluginServer {
	return &server{impl: impl}
}

func (s *server) GetPluginInfo(_ context.Context, _ *empty.Empty) (*core.PluginInfo, error) {
	return s.impl.PluginInfo(), nil
}

func (s *server) InitPlugin(_ context.Context, _ *empty.Empty) (*core.InitPluginResponse, error) {
	config, err := s.impl.Init()
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return &core.InitPluginResponse{Config: config}, nil
}

func (s *server) ResetPlugin(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	s.impl.Cleanup()

	return &empty.Empty{}, nil
}

func (s *server) ListBackdrops(_ context.Context, _ *empty.Empty) (*configuration.ListBackdropsResponse, error) {
	backdrops, err := s.impl.ListBackdrops()
	if err != nil {
		return nil, fmt.Errorf("could not list backdrops: %w", err)
	}

	return &configuration.ListBackdropsResponse{Backdrops: backdrops}, nil
}

func (s *server) GetBackdrop(_ context.Context, request *configuration.GetBackdropRequest) (*core.Backdrop, error) {
	response, err := s.impl.GetBackdrop(request.GetAlias())
	if err != nil {
		return nil, fmt.Errorf("could not get backdrop: %w", err)
	}

	return response, nil
}
