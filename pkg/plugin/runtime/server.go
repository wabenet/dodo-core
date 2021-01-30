package runtime

import (
	"net"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	"golang.org/x/net/context"
)

const (
	streamListenAddress = "127.0.0.1:"

	ErrNoStreamingConnection StreamingError = "no streaming connection established"
)

type StreamingError string

func (e StreamingError) Error() string {
	return string(e)
}

type server struct {
	impl             ContainerRuntime
	streamListener   net.Listener
	streamConnection net.Conn
}

func (s *server) Init(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.impl.Init()
}

func (s *server) GetPluginInfo(_ context.Context, _ *empty.Empty) (*api.PluginInfo, error) {
	return s.impl.PluginInfo()
}

func (s *server) GetImage(_ context.Context, request *api.GetImageRequest) (*api.GetImageResponse, error) {
	id, err := s.impl.ResolveImage(request.ImageSpec)
	if err != nil {
		return nil, err
	}

	return &api.GetImageResponse{ImageId: id}, nil
}

func (s *server) CreateContainer(_ context.Context, config *api.CreateContainerRequest) (*api.CreateContainerResponse, error) {
	id, err := s.impl.CreateContainer(config.Config, config.Tty, config.Stdio)
	if err != nil {
		return nil, err
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

func (s *server) GetStreamingConnection(_ context.Context, request *api.GetStreamingConnectionRequest) (*api.GetStreamingConnectionResponse, error) {
	listener, err := net.Listen("tcp", streamListenAddress)
	if err != nil {
		return nil, err
	}

	s.streamListener = listener

	go func() {
		conn, err := s.streamListener.Accept()
		if err != nil {
			log.Default().Error("could not accept client connection", "error", err)
		}

		s.streamConnection = conn
	}()

	return &api.GetStreamingConnectionResponse{Url: s.streamListener.Addr().String()}, nil
}

func (s *server) StreamContainer(_ context.Context, request *api.StreamContainerRequest) (*api.StreamContainerResponse, error) {
	if s.streamConnection == nil {
		return nil, ErrNoStreamingConnection
	}

	defer func() {
		if err := s.streamConnection.Close(); err != nil {
			log.Default().Error("could not close client connection", "error", err)
		}

		if err := s.streamListener.Close(); err != nil {
			log.Default().Error("could not close listener", "error", err)
		}
	}()

	err := s.impl.StreamContainer(request.ContainerId, s.streamConnection, s.streamConnection, request.Height, request.Width)

	if result, ok := err.(Result); ok {
		return &api.StreamContainerResponse{
			ExitCode: result.ExitCode,
			Message:  result.Message,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &api.StreamContainerResponse{ExitCode: 0}, nil
}
