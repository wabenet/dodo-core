package builder

import (
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/wabenet/dodo-core/api/v1alpha4"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
)

type server struct {
	impl   ImageBuilder
	stdout *grpcutil.StreamOutputServer
}

func NewGRPCServer(impl ImageBuilder) api.BuilderPluginServer {
	return &server{impl: impl}
}

func (s *server) reset() {
	s.stdout = grpcutil.NewStreamOutputServer()
}

func (s *server) GetPluginInfo(_ context.Context, _ *empty.Empty) (*api.PluginInfo, error) {
	return s.impl.PluginInfo(), nil
}

func (s *server) InitPlugin(_ context.Context, _ *empty.Empty) (*api.InitPluginResponse, error) {
	s.reset()

	config, err := s.impl.Init()
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return &api.InitPluginResponse{Config: config}, nil
}

func (s *server) ResetPlugin(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	s.reset()
	s.impl.Cleanup()

	return &empty.Empty{}, nil
}

func (s *server) StreamOutput(_ *empty.Empty, srv api.BuilderPlugin_StreamOutputServer) error {
	if err := s.stdout.SendTo(srv); err != nil {
		return fmt.Errorf("error during output stream: %w", err)
	}

	return nil
}

func (s *server) CreateImage(_ context.Context, request *api.CreateImageRequest) (*api.CreateImageResponse, error) {
	resp := &api.CreateImageResponse{}

	if request.Height == 0 && request.Width == 0 {
		id, err := s.impl.CreateImage(request.Config, nil)
		if err != nil {
			return nil, fmt.Errorf("could not build image: %w", err)
		}

		resp.ImageId = id

		return resp, nil
	}

	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		if err := s.stdout.ReadFrom(outReader, errReader); err != nil {
			return fmt.Errorf("error reading output stream: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		defer outWriter.Close()
		defer errWriter.Close()

		id, err := s.impl.CreateImage(request.Config, &plugin.StreamConfig{
			Stdout:         outWriter,
			Stderr:         errWriter,
			TerminalHeight: request.Height,
			TerminalWidth:  request.Width,
		})
		if err != nil {
			return fmt.Errorf("could not build image: %w", err)
		}

		resp.ImageId = id

		return nil
	})

	if err := eg.Wait(); err != nil {
		return resp, fmt.Errorf("error during image build stream: %w", err)
	}

	return resp, nil
}
