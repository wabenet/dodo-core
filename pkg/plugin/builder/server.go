package builder

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/wabenet/dodo-core/api/v1alpha3"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
)

type server struct {
	impl     ImageBuilder
	stdoutCh chan []byte
	stderrCh chan []byte
}

func NewGRPCServer(impl ImageBuilder) api.BuilderPluginServer {
	return &server{
		impl:     impl,
		stdoutCh: make(chan []byte),
		stderrCh: make(chan []byte),
	}
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

func (s *server) StreamBuildOutput(_ *empty.Empty, srv api.BuilderPlugin_StreamBuildOutputServer) error {
	var data api.OutputData

	for s.stdoutCh != nil && s.stderrCh != nil {
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

	return nil
}

func (s *server) CreateImage(_ context.Context, request *api.CreateImageRequest) (*api.CreateImageResponse, error) {
	resp := &api.CreateImageResponse{}

	if request.Height == 0 && request.Width == 0 {
		if id, err := s.impl.CreateImage(request.Config, nil); err != nil {
			return nil, fmt.Errorf("could not build image: %w", err)
		} else {
			resp.ImageId = id
		}

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
		defer outWriter.Close()
		defer errWriter.Close()

		if id, err := s.impl.CreateImage(request.Config, &plugin.StreamConfig{
			Stdout:         outWriter,
			Stderr:         errWriter,
			TerminalHeight: request.Height,
			TerminalWidth:  request.Width,
		}); err != nil {
			return fmt.Errorf("could not build image: %w", err)
		} else {
			resp.ImageId = id
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return resp, err
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
			return err
		}
	}
}
