package runtime

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/wabenet/dodo-core/api/v1alpha4"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
)

var ErrUnexpectedMapType = errors.New("unexpected map type for stdio streaming server")

type server struct {
	impl   ContainerRuntime
	stdin  sync.Map
	stdout sync.Map
}

func NewGRPCServer(impl ContainerRuntime) api.RuntimePluginServer {
	return &server{impl: impl}
}

func (s *server) reset() {
	s.stdin = sync.Map{}
	s.stdout = sync.Map{}
}

func (s *server) stdinServer(containerID string) (*grpcutil.StreamInputServer, error) {
	inputServer, _ := s.stdin.LoadOrStore(containerID, grpcutil.NewStreamInputServer())

	result, ok := inputServer.(*grpcutil.StreamInputServer)
	if !ok {
		return nil, ErrUnexpectedMapType
	}

	return result, nil
}

func (s *server) stdoutServer(containerID string) (*grpcutil.StreamOutputServer, error) {
	outputServer, _ := s.stdout.LoadOrStore(containerID, grpcutil.NewStreamOutputServer())

	result, ok := outputServer.(*grpcutil.StreamOutputServer)
	if !ok {
		return nil, ErrUnexpectedMapType
	}

	return result, nil
}

type streamInputServer struct {
	server api.RuntimePlugin_StreamInputServer
}

func (s *streamInputServer) Recv() (*api.InputData, error) {
	d, err := s.server.Recv()
	if err != nil {
		return nil, fmt.Errorf("error wrapping Recv call: %w", err)
	}

	return d.GetInputData(), nil
}

func (s *streamInputServer) SendAndClose(e *empty.Empty) error {
	if err := s.server.SendAndClose(e); err != nil {
		return fmt.Errorf("error wrapping SendAndClose call: %w", err)
	}

	return nil
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

func (s *server) StreamInput(srv api.RuntimePlugin_StreamInputServer) error {
	req, err := srv.Recv()
	if err != nil {
		return fmt.Errorf("error during input stream: %w", err)
	}

	id := req.GetInitialRequest().Id

	inputServer, err := s.stdinServer(id)
	if err != nil {
		return fmt.Errorf("could not find stream input server: %w", err)
	}

	if err := inputServer.ReceiveFrom(&streamInputServer{server: srv}); err != nil {
		return fmt.Errorf("error during input stream: %w", err)
	}

	return nil
}

func (s *server) StreamOutput(request *api.StreamOutputRequest, srv api.RuntimePlugin_StreamOutputServer) error {
	id := request.Id

	outputServer, err := s.stdoutServer(id)
	if err != nil {
		return fmt.Errorf("could not find stream output server: %w", err)
	}

	if err := outputServer.SendTo(srv); err != nil {
		return fmt.Errorf("error during output stream: %w", err)
	}

	return nil
}

func (s *server) StreamContainer(
	_ context.Context,
	request *api.StreamContainerRequest,
) (*api.StreamContainerResponse, error) {
	resp := &api.StreamContainerResponse{}

	inReader, inWriter := io.Pipe()
	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()

	inputServer, err := s.stdinServer(request.ContainerId)
	if err != nil {
		return nil, fmt.Errorf("could not find stream input server: %w", err)
	}

	outputServer, err := s.stdoutServer(request.ContainerId)
	if err != nil {
		return nil, fmt.Errorf("could not find stream output server: %w", err)
	}

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error { return copyInputServerToStdin(inputServer, inWriter) })
	eg.Go(func() error { return copyOutputServerToStdout(outputServer, outReader, errReader) })

	eg.Go(func() error {
		defer outWriter.Close()
		defer errWriter.Close()
		defer inputServer.Close()

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

func copyInputServerToStdin(inputServer *grpcutil.StreamInputServer, stdin io.WriteCloser) error {
	if err := inputServer.WriteTo(stdin); err != nil {
		return fmt.Errorf("error writing input stream: %w", err)
	}

	return nil
}

func copyOutputServerToStdout(outputServer *grpcutil.StreamOutputServer, stdout, stderr io.Reader) error {
	if err := outputServer.ReadFrom(stdout, stderr); err != nil {
		return fmt.Errorf("error reading output stream: %w", err)
	}

	return nil
}
