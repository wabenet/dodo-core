package configuration

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/wabenet/dodo-core/api/configuration/v1alpha2"
	pluginapi "github.com/wabenet/dodo-core/api/plugin/v1alpha1"
)

type Server struct {
	impl Configuration
}

func NewGRPCServer(impl Configuration) *Server {
	return &Server{impl: impl}
}

func (s *Server) GetPluginInfo(_ context.Context, _ *empty.Empty) (*pluginapi.PluginInfo, error) {
	return s.impl.PluginInfo(), nil
}

func (s *Server) InitPlugin(_ context.Context, _ *empty.Empty) (*pluginapi.InitPluginResponse, error) {
	config, err := s.impl.Init()
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return &pluginapi.InitPluginResponse{Config: config}, nil
}

func (s *Server) ResetPlugin(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	s.impl.Cleanup()

	return &empty.Empty{}, nil
}

func (s *Server) ListBackdrops(_ context.Context, _ *empty.Empty) (*api.ListBackdropsResponse, error) {
	backdrops, err := s.impl.ListBackdrops()
	if err != nil {
		return nil, fmt.Errorf("could not list backdrops: %w", err)
	}

	result := []*api.Backdrop{}

	for _, b := range backdrops {
		result = append(result, b.ToProto())
	}

	return &api.ListBackdropsResponse{Backdrops: result}, nil
}

func (s *Server) GetBackdrop(_ context.Context, request *api.GetBackdropRequest) (*api.GetBackdropResponse, error) {
	backdrop, err := s.impl.GetBackdrop(request.GetAlias())
	if err != nil {
		return nil, fmt.Errorf("could not get backdrop: %w", err)
	}

	return &api.GetBackdropResponse{Backdrop: backdrop.ToProto()}, nil
}
