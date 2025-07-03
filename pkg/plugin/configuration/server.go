package configuration

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/configuration/v1alpha2"
	pluginapi "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
)

type Server struct {
	api.UnsafePluginServer

	impl Configuration
}

func NewGRPCServer(impl Configuration) *Server {
	return &Server{impl: impl}
}

func (s *Server) GetPluginMetadata(_ context.Context, _ *empty.Empty) (*pluginapi.PluginMetadata, error) {
	return s.impl.Metadata().ToProto(), nil
}

func (s *Server) InitPlugin(_ context.Context, _ *empty.Empty) (*pluginapi.InitPluginResponse, error) {
	config, err := s.impl.Init()
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	resp := &pluginapi.InitPluginResponse{}

	resp.SetConfig(config)

	return resp, nil
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

	resp := &api.ListBackdropsResponse{}

	resp.SetBackdrops(result)

	return resp, nil
}

func (s *Server) GetBackdrop(_ context.Context, request *api.GetBackdropRequest) (*api.GetBackdropResponse, error) {
	backdrop, err := s.impl.GetBackdrop(request.GetAlias())
	if err != nil {
		return nil, fmt.Errorf("could not get backdrop: %w", err)
	}

	resp := &api.GetBackdropResponse{}

	resp.SetBackdrop(backdrop.ToProto())

	return resp, nil
}
