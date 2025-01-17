package runtime

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	pluginapi "github.com/wabenet/dodo-core/api/plugin/v1alpha1"
	runtime "github.com/wabenet/dodo-core/api/runtime/v1alpha2"
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

func NewGRPCServer(impl ContainerRuntime) runtime.PluginServer {
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
	server runtime.Plugin_StreamInputServer
}

func (s *streamInputServer) Recv() (*pluginapi.InputData, error) {
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

func (s *server) GetPluginInfo(_ context.Context, _ *empty.Empty) (*pluginapi.PluginInfo, error) {
	return s.impl.PluginInfo(), nil
}

func (s *server) InitPlugin(_ context.Context, _ *empty.Empty) (*pluginapi.InitPluginResponse, error) {
	s.reset()

	config, err := s.impl.Init()
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return &pluginapi.InitPluginResponse{Config: config}, nil
}

func (s *server) ResetPlugin(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	s.reset()
	s.impl.Cleanup()

	return &empty.Empty{}, nil
}

func (s *server) GetImage(_ context.Context, request *runtime.GetImageRequest) (*runtime.GetImageResponse, error) {
	id, err := s.impl.ResolveImage(request.GetImageSpec())
	if err != nil {
		return nil, fmt.Errorf("could not resolve image: %w", err)
	}

	return &runtime.GetImageResponse{ImageId: id}, nil
}

func (s *server) CreateContainer(
	_ context.Context,
	config *runtime.CreateContainerRequest,
) (*runtime.CreateContainerResponse, error) {
	id, err := s.impl.CreateContainer(config.GetConfig(), config.GetTty(), config.GetStdio())
	if err != nil {
		return nil, fmt.Errorf("could not create container: %w", err)
	}

	return &runtime.CreateContainerResponse{ContainerId: id}, nil
}

func (s *server) StartContainer(_ context.Context, request *runtime.StartContainerRequest) (*empty.Empty, error) {
	if err := s.impl.StartContainer(request.GetContainerId()); err != nil {
		return nil, fmt.Errorf("could not start container: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *server) DeleteContainer(_ context.Context, request *runtime.DeleteContainerRequest) (*empty.Empty, error) {
	if err := s.impl.DeleteContainer(request.GetContainerId()); err != nil {
		return nil, fmt.Errorf("could not delete container: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *server) ResizeContainer(_ context.Context, request *runtime.ResizeContainerRequest) (*empty.Empty, error) {
	if err := s.impl.ResizeContainer(request.GetContainerId(), request.GetHeight(), request.GetWidth()); err != nil {
		return nil, fmt.Errorf("could not resize container: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *server) KillContainer(_ context.Context, request *runtime.KillContainerRequest) (*empty.Empty, error) {
	if err := s.impl.KillContainer(request.GetContainerId(), signalFromString(request.GetSignal())); err != nil {
		return nil, fmt.Errorf("could not kill container: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *server) StreamInput(srv runtime.Plugin_StreamInputServer) error {
	req, err := srv.Recv()
	if err != nil {
		return fmt.Errorf("error during input stream: %w", err)
	}

	id := req.GetInitialRequest().GetId()

	inputServer, err := s.stdinServer(id)
	if err != nil {
		return fmt.Errorf("could not find stream input server: %w", err)
	}

	if err := inputServer.ReceiveFrom(&streamInputServer{server: srv}); err != nil {
		return fmt.Errorf("error during input stream: %w", err)
	}

	return nil
}

func (s *server) StreamOutput(request *pluginapi.StreamOutputRequest, srv runtime.Plugin_StreamOutputServer) error {
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

func (s *server) StreamContainer(
	_ context.Context,
	request *runtime.StreamContainerRequest,
) (*runtime.StreamContainerResponse, error) {
	resp := &runtime.StreamContainerResponse{}

	inReader, inWriter := io.Pipe()
	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()

	inputServer, err := s.stdinServer(request.GetContainerId())
	if err != nil {
		return nil, fmt.Errorf("could not find stream input server: %w", err)
	}

	outputServer, err := s.stdoutServer(request.GetContainerId())
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

		streamResp, err := s.impl.StreamContainer(request.GetContainerId(), &plugin.StreamConfig{
			Stdin:          inReader,
			Stdout:         outWriter,
			Stderr:         errWriter,
			TerminalHeight: request.GetHeight(),
			TerminalWidth:  request.GetWidth(),
		})
		if err != nil {
			return fmt.Errorf("could not stream container: %w", err)
		}

		resp.ExitCode = int64(streamResp.ExitCode)

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

func (s *server) CreateVolume(_ context.Context, request *runtime.CreateVolumeRequest) (*empty.Empty, error) {
	if err := s.impl.CreateVolume(request.GetName()); err != nil {
		return nil, fmt.Errorf("could create volume: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *server) DeleteVolume(_ context.Context, request *runtime.DeleteVolumeRequest) (*empty.Empty, error) {
	if err := s.impl.CreateVolume(request.GetName()); err != nil {
		return nil, fmt.Errorf("could delete volume: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *server) WriteFile(_ context.Context, request *runtime.WriteFileRequest) (*empty.Empty, error) {
	if err := s.impl.WriteFile(
		request.GetContainerId(),
		request.GetFilePath(),
		[]byte(request.GetContents()),
	); err != nil {
		return nil, fmt.Errorf("could not write file: %w", err)
	}

	return &empty.Empty{}, nil
}
