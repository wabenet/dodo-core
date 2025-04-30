package builder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/wabenet/dodo-core/api/build/v1alpha2"
	pluginapi "github.com/wabenet/dodo-core/api/plugin/v1alpha1"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
)

var ErrUnexpectedMapType = errors.New("unexpected map type for stdio streaming server")

type Server struct {
	impl   ImageBuilder
	stdout sync.Map
}

func NewGRPCServer(impl ImageBuilder) *Server {
	return &Server{impl: impl}
}

func (s *Server) reset() {
	s.stdout = sync.Map{}
}

func (s *Server) stdoutServer(streamID string) (*grpcutil.StreamOutputServer, error) {
	outputServer, _ := s.stdout.LoadOrStore(streamID, grpcutil.NewStreamOutputServer())

	result, ok := outputServer.(*grpcutil.StreamOutputServer)
	if !ok {
		return nil, ErrUnexpectedMapType
	}

	return result, nil
}

func (s *Server) GetPluginInfo(_ context.Context, _ *empty.Empty) (*pluginapi.PluginInfo, error) {
	return s.impl.PluginInfo(), nil
}

func (s *Server) InitPlugin(_ context.Context, _ *empty.Empty) (*pluginapi.InitPluginResponse, error) {
	s.reset()

	config, err := s.impl.Init()
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return &pluginapi.InitPluginResponse{Config: config}, nil
}

func (s *Server) ResetPlugin(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	s.reset()
	s.impl.Cleanup()

	return &empty.Empty{}, nil
}

func (s *Server) StreamOutput(request *pluginapi.StreamOutputRequest, srv api.Plugin_StreamOutputServer) error {
	id := request.GetId()

	outputServer, err := s.stdoutServer(id)
	if err != nil {
		return fmt.Errorf("could not find stream output server: %w", err)
	}

	if err := outputServer.SendTo(srv); err != nil {
		return fmt.Errorf("error during output stream: %w", err)
	}

	return nil
}

func (s *Server) CreateImage(_ context.Context, request *api.CreateImageRequest) (*api.CreateImageResponse, error) {
	resp := &api.CreateImageResponse{}

	if request.GetHeight() == 0 && request.GetWidth() == 0 {
		id, err := s.impl.CreateImage(BuildConfigFromProto(request.GetConfig()), nil)
		if err != nil {
			return nil, fmt.Errorf("could not build image: %w", err)
		}

		resp.ImageId = id

		return resp, nil
	}

	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()

	outputServer, err := s.stdoutServer(request.GetStreamId())
	if err != nil {
		return nil, fmt.Errorf("could not find stream output server: %w", err)
	}

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error { return copyOutputServerToStdout(outputServer, outReader, errReader) })

	eg.Go(func() error {
		defer outWriter.Close()
		defer errWriter.Close()

		imageID, err := s.impl.CreateImage(BuildConfigFromProto(request.GetConfig()), &plugin.StreamConfig{
			Stdout:         outWriter,
			Stderr:         errWriter,
			TerminalHeight: request.GetHeight(),
			TerminalWidth:  request.GetWidth(),
		})
		if err != nil {
			return fmt.Errorf("could not build image: %w", err)
		}

		resp.ImageId = imageID

		return nil
	})

	if err := eg.Wait(); err != nil {
		return resp, fmt.Errorf("error during image build stream: %w", err)
	}

	return resp, nil
}

func copyOutputServerToStdout(outputServer *grpcutil.StreamOutputServer, stdout, stderr io.Reader) error {
	if err := outputServer.ReadFrom(stdout, stderr); err != nil {
		return fmt.Errorf("error reading output stream: %w", err)
	}

	return nil
}
