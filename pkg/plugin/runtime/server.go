package runtime

import (
	"bufio"
	"context"
	"fmt"
	"io"

	api "github.com/dodo-cli/dodo-core/api/v1alpha3"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	impl     ContainerRuntime
	stdinCh  chan []byte
	stdoutCh chan []byte
	stderrCh chan []byte
}

func NewGRPCServer(impl ContainerRuntime, listen string) api.RuntimePluginServer {
	return &server{
		impl:     impl,
		stdinCh:  make(chan []byte),
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

func (s *server) GetImage(_ context.Context, request *api.GetImageRequest) (*api.GetImageResponse, error) {
	id, err := s.impl.ResolveImage(request.ImageSpec)
	if err != nil {
		return nil, fmt.Errorf("could not resolve image: %w", err)
	}

	return &api.GetImageResponse{ImageId: id}, nil
}

func (s *server) CreateContainer(_ context.Context, config *api.CreateContainerRequest) (*api.CreateContainerResponse, error) {
	id, err := s.impl.CreateContainer(config.Config, config.Tty, config.Stdio)
	if err != nil {
		return nil, fmt.Errorf("could not create container: %w", err)
	}

	return &api.CreateContainerResponse{ContainerId: id}, nil
}

func (s *server) StartContainer(_ context.Context, request *api.StartContainerRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.impl.StartContainer(request.ContainerId)
}

func (s *server) DeleteContainer(_ context.Context, request *api.DeleteContainerRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.impl.DeleteContainer(request.ContainerId)
}

func (s *server) ResizeContainer(_ context.Context, request *api.ResizeContainerRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.impl.ResizeContainer(request.ContainerId, request.Height, request.Width)
}

func (s *server) KillContainer(_ context.Context, request *api.KillContainerRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.impl.KillContainer(request.ContainerId, signalFromString(request.Signal))
}

func (s *server) StreamRuntimeInput(srv api.RuntimePlugin_StreamRuntimeInputServer) error {
	defer close(s.stdinCh)

	for {
		data, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				return srv.SendAndClose(&empty.Empty{})

			}

			if status.Code(err) == codes.Unavailable ||
				status.Code(err) == codes.Canceled ||
				status.Code(err) == codes.Unimplemented ||
				err == context.Canceled {
				return nil
			}

			log.L().Error("error receiving data", "err", err)
			return err
		}

		s.stdinCh <- data.Data
	}
}

func (s *server) StreamRuntimeOutput(_ *empty.Empty, srv api.RuntimePlugin_StreamRuntimeOutputServer) error {
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
			return err
		}
	}
}

func (s *server) StreamContainer(_ context.Context, request *api.StreamContainerRequest) (*api.StreamContainerResponse, error) {
	inReader, inWriter := io.Pipe()
	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()

	go func() {
		copyInput(inWriter, s.stdinCh)
		inWriter.Close()
	}()

	go copyOutput(s.stdoutCh, outReader)
	go copyOutput(s.stderrCh, errReader)

	defer func() {
		outWriter.Close()
		errWriter.Close()
	}()

	r, err := s.impl.StreamContainer(request.ContainerId, &plugin.StreamConfig{
		Stdin:          inReader,
		Stdout:         outWriter,
		Stderr:         errWriter,
		TerminalHeight: request.Height,
		TerminalWidth:  request.Width,
	})
	if err != nil {
		return nil, fmt.Errorf("could not stream container: %w", err)
	}

	return &api.StreamContainerResponse{ExitCode: int64(r.ExitCode)}, nil
}

func copyInput(dst io.Writer, src chan []byte) {
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

	bufdst.Flush()
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
