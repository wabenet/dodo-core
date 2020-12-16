package runtime

import (
	"net"

	"github.com/dodo-cli/dodo-core/pkg/types"
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

func (s *server) Init(_ context.Context, _ *types.Empty) (*types.PluginInfo, error) {
	return s.impl.Init()
}

func (s *server) ResolveImage(_ context.Context, request *types.Image) (*types.Image, error) {
	id, err := s.impl.ResolveImage(request.Name)
	if err != nil {
		return nil, err
	}

	return &types.Image{Name: request.Name, Id: id}, nil
}

func (s *server) CreateContainer(_ context.Context, config *types.ContainerConfig) (*types.ContainerId, error) {
	id, err := s.impl.CreateContainer(config.Config, config.Tty, config.Stdio)
	if err != nil {
		return nil, err
	}

	return &types.ContainerId{Id: id}, nil
}

func (s *server) StartContainer(_ context.Context, request *types.ContainerId) (*types.Empty, error) {
	return &types.Empty{}, s.impl.StartContainer(request.Id)
}

func (s *server) RemoveContainer(_ context.Context, request *types.ContainerId) (*types.Empty, error) {
	return &types.Empty{}, s.impl.RemoveContainer(request.Id)
}

func (s *server) ResizeContainer(_ context.Context, request *types.ContainerBox) (*types.Empty, error) {
	return &types.Empty{}, s.impl.ResizeContainer(request.Id, request.Height, request.Width)
}

func (s *server) SetupStreamingConnection(_ context.Context, request *types.ContainerId) (*types.StreamingConnection, error) {
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

	return &types.StreamingConnection{Url: s.streamListener.Addr().String()}, nil
}

func (s *server) StreamContainer(_ context.Context, request *types.ContainerBox) (*types.Result, error) {
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

	err := s.impl.StreamContainer(request.Id, s.streamConnection, s.streamConnection, request.Height, request.Width)

	if result, ok := err.(types.Result); ok {
		return &result, nil
	}

	if err != nil {
		return nil, err
	}

	return &types.Result{ExitCode: 0}, nil
}
