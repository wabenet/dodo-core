package builder

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/build/v1alpha2"
	pluginapi "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/stdio"
	"golang.org/x/sync/errgroup"
)

var ErrUnexpectedMapType = errors.New("unexpected map type for stdio streaming server")

type Server struct {
	pluginapi.UnsafePluginServer
	pluginapi.UnsafeOutputStreamingPluginServer
	api.UnsafeBuilderPluginServer

	stdio.OutputStreamingServer

	impl ImageBuilder
}

func NewGRPCServer(impl ImageBuilder) *Server {
	return &Server{impl: impl}
}

func (s *Server) reset() {
	s.OutputStreamingServer.Reset()
}

func (s *Server) GetPluginMetadata(_ context.Context, _ *empty.Empty) (*pluginapi.GetPluginMetadataResponse, error) {
	resp := &pluginapi.GetPluginMetadataResponse{}

	resp.SetMetadata(s.impl.Metadata().ToProto())

	return resp, nil
}

func (s *Server) InitPlugin(_ context.Context, _ *empty.Empty) (*pluginapi.InitPluginResponse, error) {
	s.reset()

	config, err := s.impl.Init()
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	pluginConfig := &pluginapi.PluginConfig{}
	resp := &pluginapi.InitPluginResponse{}

	pluginConfig.SetConfig(config)
	resp.SetConfig(pluginConfig)

	return resp, nil
}

func (s *Server) ResetPlugin(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	s.reset()
	s.impl.Cleanup()

	return &empty.Empty{}, nil
}

func (s *Server) CreateImage(_ context.Context, request *api.CreateImageRequest) (*api.CreateImageResponse, error) {
	resp := &api.CreateImageResponse{}

	if request.GetHeight() == 0 && request.GetWidth() == 0 {
		id, err := s.impl.CreateImage(BuildConfigFromProto(request.GetConfig()), nil)
		if err != nil {
			return nil, fmt.Errorf("could not build image: %w", err)
		}

		resp.SetImageId(id)

		return resp, nil
	}

	outputStream, err := s.OutputStreamingServer.PrepareStream(request.GetStreamId())
	if err != nil {
		return nil, err
	}

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(outputStream.Copy)

	eg.Go(func() error {
		defer outputStream.Close()

		imageID, err := s.impl.CreateImage(BuildConfigFromProto(request.GetConfig()), &plugin.StreamConfig{
			Stdout:         outputStream.Stdout,
			Stderr:         outputStream.Stderr,
			TerminalHeight: request.GetHeight(),
			TerminalWidth:  request.GetWidth(),
		})
		if err != nil {
			return fmt.Errorf("could not build image: %w", err)
		}

		resp.SetImageId(imageID)

		return nil
	})

	if err := eg.Wait(); err != nil {
		return resp, fmt.Errorf("error during image build stream: %w", err)
	}

	return resp, nil
}
