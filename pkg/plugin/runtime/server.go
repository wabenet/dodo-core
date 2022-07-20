package runtime

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	api "github.com/wabenet/dodo-core/api/v1alpha4"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	impl        ContainerRuntime
	stdinCh     chan []byte
	stdoutCh    chan []byte
	stderrCh    chan []byte
	inputDone   chan error
	outputDone  chan error
	stdinCloser sync.Once
}

func NewGRPCServer(impl ContainerRuntime) api.RuntimePluginServer {
	return &server{impl: impl}
}

func (s *server) reset() {
	s.stdinCh = make(chan []byte)
	s.stdoutCh = make(chan []byte)
	s.stderrCh = make(chan []byte)
	s.inputDone = make(chan error, 1)
	s.outputDone = make(chan error, 1)
	s.stdinCloser = sync.Once{}
}

func (s *server) closeStdin() {
	s.stdinCloser.Do(func() {
		close(s.stdinCh)
		s.inputDone <- nil
	})
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

func (s *server) GetImage(_ context.Context, request *api.GetImageRequest) (*api.GetImageResponse, error) {
	id, err := s.impl.ResolveImage(request.ImageSpec)
	if err != nil {
		return nil, fmt.Errorf("could not resolve image: %w", err)
	}

	return &api.GetImageResponse{ImageId: id}, nil
}

func (s *server) CreateContainer(
	_ context.Context,
	config *api.CreateContainerRequest,
) (*api.CreateContainerResponse, error) {
	id, err := s.impl.CreateContainer(config.Config, config.Tty, config.Stdio)
	if err != nil {
		return nil, fmt.Errorf("could not create container: %w", err)
	}

	return &api.CreateContainerResponse{ContainerId: id}, nil
}

func (s *server) StartContainer(_ context.Context, request *api.StartContainerRequest) (*empty.Empty, error) {
	if err := s.impl.StartContainer(request.ContainerId); err != nil {
		return nil, fmt.Errorf("could not start container: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *server) DeleteContainer(_ context.Context, request *api.DeleteContainerRequest) (*empty.Empty, error) {
	if err := s.impl.DeleteContainer(request.ContainerId); err != nil {
		return nil, fmt.Errorf("could not delete container: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *server) ResizeContainer(_ context.Context, request *api.ResizeContainerRequest) (*empty.Empty, error) {
	if err := s.impl.ResizeContainer(request.ContainerId, request.Height, request.Width); err != nil {
		return nil, fmt.Errorf("could not resize container: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *server) KillContainer(_ context.Context, request *api.KillContainerRequest) (*empty.Empty, error) {
	if err := s.impl.KillContainer(request.ContainerId, signalFromString(request.Signal)); err != nil {
		return nil, fmt.Errorf("could not kill container: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *server) StreamRuntimeInput(srv api.RuntimePlugin_StreamRuntimeInputServer) error {
	defer s.closeStdin()

	for {
		data, err := srv.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if err := srv.SendAndClose(&empty.Empty{}); err != nil {
					return fmt.Errorf("could not close input stream: %w", err)
				}

				return nil
			}

			if errors.Is(err, context.Canceled) ||
				status.Code(err) == codes.Unavailable ||
				status.Code(err) == codes.Canceled ||
				status.Code(err) == codes.Unimplemented {
				return nil
			}

			log.L().Error("error receiving data", "err", err)

			return fmt.Errorf("error receiving build input from clien: %w", err)
		}

		s.stdinCh <- data.Data
	}
}

func (s *server) StreamRuntimeOutput(_ *empty.Empty, srv api.RuntimePlugin_StreamRuntimeOutputServer) error {
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

func (s *server) StreamContainer(
	_ context.Context,
	request *api.StreamContainerRequest,
) (*api.StreamContainerResponse, error) {
	resp := &api.StreamContainerResponse{}

	inReader, inWriter := io.Pipe()
	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		defer inWriter.Close()

		return copyInput(inWriter, s.stdinCh)
	})

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
		return <-s.inputDone
	})

	eg.Go(func() error {
		defer outWriter.Close()
		defer errWriter.Close()

		defer s.closeStdin()

		r, err := s.impl.StreamContainer(request.ContainerId, &plugin.StreamConfig{
			Stdin:          inReader,
			Stdout:         outWriter,
			Stderr:         errWriter,
			TerminalHeight: request.Height,
			TerminalWidth:  request.Width,
		})
		if err != nil {
			return fmt.Errorf("could not stream container: %w", err)
		}

		resp.ExitCode = int64(r.ExitCode)

		return nil
	})

	if err := eg.Wait(); err != nil {
		return resp, fmt.Errorf("error during container stream: %w", err)
	}

	return resp, nil
}

func copyInput(dst io.Writer, src chan []byte) error {
	bufdst := bufio.NewWriter(dst)

	for data := range src {
		if len(data) == 0 {
			continue
		}

		if _, err := bufdst.Write(data); err != nil {
			log.L().Warn("error in stdio stream", "err", err)

			break
		}
	}

	if err := bufdst.Flush(); err != nil {
		return fmt.Errorf("could not flush container input: %w", err)
	}

	return nil
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
			return fmt.Errorf("error copying container output: %w", err)
		}
	}
}
