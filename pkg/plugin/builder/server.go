package builder

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/wabenet/dodo-core/api/v1alpha4"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
)

type server struct {
	impl       ImageBuilder
	stdoutCh   chan []byte
	stderrCh   chan []byte
	outputDone chan error
}

func NewGRPCServer(impl ImageBuilder) api.BuilderPluginServer {
	return &server{impl: impl}
}

func (s *server) reset() {
	s.stdoutCh = make(chan []byte)
	s.stderrCh = make(chan []byte)
	s.outputDone = make(chan error, 1)
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

func (s *server) StreamBuildOutput(_ *empty.Empty, srv api.BuilderPlugin_StreamBuildOutputServer) error {
	var data api.OutputData

	defer func() {
		s.outputDone <- nil
	}()

	for {
		if s.stdoutCh == nil && s.stderrCh == nil {
			return nil
		}

		select {
		case d, ok := <-s.stdoutCh:
			if !ok {
				s.stdoutCh = nil

				continue
			}

			data.Data = d
			data.Channel = api.OutputData_STDOUT

		case d, ok := <-s.stderrCh:
			if !ok {
				s.stderrCh = nil

				continue
			}

			data.Data = d
			data.Channel = api.OutputData_STDERR

		case <-srv.Context().Done():
			return nil
		}

		if len(data.Data) == 0 {
			continue
		}

		if err := srv.Send(&data); err != nil {
			return fmt.Errorf("error sending build output to client: %w", err)
		}
	}
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
		return copyOutput(s.stdoutCh, outReader)
	})

	eg.Go(func() error {
		return copyOutput(s.stderrCh, errReader)
	})

	eg.Go(func() error {
		return <-s.outputDone
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

func copyOutput(dst chan []byte, src io.Reader) error {
	defer close(dst)

	bufsrc := bufio.NewReader(src)

	for {
		var data [1024]byte

		n, err := bufsrc.Read(data[:])

		if n > 0 {
			dst <- data[:n]
		}

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return fmt.Errorf("error copying build output: %w", err)
		}
	}
}
