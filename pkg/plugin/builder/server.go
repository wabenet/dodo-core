package builder

import (
	"bufio"
	"context"
	"fmt"
	"io"

	api "github.com/dodo-cli/dodo-core/api/v1alpha3"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
)

type server struct {
	impl     ImageBuilder
	stdoutCh chan []byte
	stderrCh chan []byte
}

func NewGRPCServer(impl ImageBuilder, listen string) api.BuilderPluginServer {
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

	for {
		select {
		case data.Data = <-s.stdoutCh:
			data.Channel = api.OutputData_STDOUT

		case data.Data = <-s.stderrCh:
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
	if request.Height == 0 && request.Width == 0 {
		imageID, err := s.impl.CreateImage(request.Config, nil)

		if err != nil {
			return nil, fmt.Errorf("could not build image: %w", err)
		}

		return &api.CreateImageResponse{ImageId: imageID}, nil
	}

	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()

	go copyOutput(s.stdoutCh, outReader)
	go copyOutput(s.stderrCh, errReader)

	defer func() {
		outWriter.Close()
		errWriter.Close()
	}()

	imageID, err := s.impl.CreateImage(request.Config, &plugin.StreamConfig{
		Stdout:         outWriter,
		Stderr:         errWriter,
		TerminalHeight: request.Height,
		TerminalWidth:  request.Width,
	})

	if err != nil {
		return nil, fmt.Errorf("could not build image: %w", err)
	}

	return &api.CreateImageResponse{ImageId: imageID}, nil
}

func copyOutput(dst chan []byte, src io.Reader) {
	bufsrc := bufio.NewReader(src)

	for {
		var data [1024]byte

		n, err := bufsrc.Read(data[:])

		if n > 0 {
			dst <- data[:n]
		}

		if err == io.EOF {
			return
		}

		if err != nil {
			log.L().Warn("error in stdio stream", "err", err)

			return
		}
	}
}
