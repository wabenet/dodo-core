package runtime

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	pluginapi "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/runtime/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
)

var ErrUnexpectedMapType = errors.New("unexpected map type for stdio streaming server")

type Server struct {
	api.UnsafePluginServer

	impl   ContainerRuntime
	stdin  sync.Map
	stdout sync.Map
}

func NewGRPCServer(impl ContainerRuntime) *Server {
	return &Server{impl: impl}
}

func (s *Server) reset() {
	s.stdin = sync.Map{}
	s.stdout = sync.Map{}
}

func (s *Server) stdinServer(containerID string) (*grpcutil.StreamInputServer, error) {
	inputServer, _ := s.stdin.LoadOrStore(containerID, grpcutil.NewStreamInputServer())

	result, ok := inputServer.(*grpcutil.StreamInputServer)
	if !ok {
		return nil, ErrUnexpectedMapType
	}

	return result, nil
}

func (s *Server) stdoutServer(containerID string) (*grpcutil.StreamOutputServer, error) {
	outputServer, _ := s.stdout.LoadOrStore(containerID, grpcutil.NewStreamOutputServer())

	result, ok := outputServer.(*grpcutil.StreamOutputServer)
	if !ok {
		return nil, ErrUnexpectedMapType
	}

	return result, nil
}

type streamInputServer struct {
	server api.Plugin_StreamInputServer
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

func (s *Server) GetPluginMetadata(_ context.Context, _ *empty.Empty) (*pluginapi.PluginMetadata, error) {
	return s.impl.Metadata().ToProto(), nil
}

func (s *Server) InitPlugin(_ context.Context, _ *empty.Empty) (*pluginapi.InitPluginResponse, error) {
	s.reset()

	config, err := s.impl.Init()
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	resp := &pluginapi.InitPluginResponse{}

	resp.SetConfig(config)

	return resp, nil
}

func (s *Server) ResetPlugin(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	s.reset()
	s.impl.Cleanup()

	return &empty.Empty{}, nil
}

func (s *Server) GetImage(_ context.Context, request *api.GetImageRequest) (*api.GetImageResponse, error) {
	id, err := s.impl.ResolveImage(request.GetImageSpec())
	if err != nil {
		return nil, fmt.Errorf("could not resolve image: %w", err)
	}

	resp := &api.GetImageResponse{}

	resp.SetImageId(id)

	return resp, nil
}

func (s *Server) CreateContainer(
	_ context.Context,
	config *api.CreateContainerRequest,
) (*api.CreateContainerResponse, error) {
	id, err := s.impl.CreateContainer(ContainerConfigFromProto(config.GetConfig()))
	if err != nil {
		return nil, fmt.Errorf("could not create container: %w", err)
	}

	resp := &api.CreateContainerResponse{}

	resp.SetContainerId(id)

	return resp, nil
}

func (s *Server) StartContainer(_ context.Context, request *api.StartContainerRequest) (*empty.Empty, error) {
	if err := s.impl.StartContainer(request.GetContainerId()); err != nil {
		return nil, fmt.Errorf("could not start container: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *Server) DeleteContainer(_ context.Context, request *api.DeleteContainerRequest) (*empty.Empty, error) {
	if err := s.impl.DeleteContainer(request.GetContainerId()); err != nil {
		return nil, fmt.Errorf("could not delete container: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *Server) ResizeContainer(_ context.Context, request *api.ResizeContainerRequest) (*empty.Empty, error) {
	if err := s.impl.ResizeContainer(request.GetContainerId(), request.GetHeight(), request.GetWidth()); err != nil {
		return nil, fmt.Errorf("could not resize container: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *Server) KillContainer(_ context.Context, request *api.KillContainerRequest) (*empty.Empty, error) {
	if err := s.impl.KillContainer(request.GetContainerId(), signalFromString(request.GetSignal())); err != nil {
		return nil, fmt.Errorf("could not kill container: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *Server) StreamInput(srv api.Plugin_StreamInputServer) error {
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

func (s *Server) StreamContainer(
	_ context.Context,
	request *api.StreamContainerRequest,
) (*api.StreamContainerResponse, error) {
	resp := &api.StreamContainerResponse{}

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

		resp.SetExitCode(int64(streamResp.ExitCode))

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

func (s *Server) CreateVolume(_ context.Context, request *api.CreateVolumeRequest) (*empty.Empty, error) {
	if err := s.impl.CreateVolume(request.GetName()); err != nil {
		return nil, fmt.Errorf("could create volume: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *Server) DeleteVolume(_ context.Context, request *api.DeleteVolumeRequest) (*empty.Empty, error) {
	if err := s.impl.CreateVolume(request.GetName()); err != nil {
		return nil, fmt.Errorf("could delete volume: %w", err)
	}

	return &empty.Empty{}, nil
}

func (s *Server) WriteFile(_ context.Context, request *api.WriteFileRequest) (*empty.Empty, error) {
	if err := s.impl.WriteFile(
		request.GetContainerId(),
		request.GetFilePath(),
		[]byte(request.GetContents()),
	); err != nil {
		return nil, fmt.Errorf("could not write file: %w", err)
	}

	return &empty.Empty{}, nil
}
