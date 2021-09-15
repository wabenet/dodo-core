package runtime

import (
	"context"
	"errors"
	"fmt"
	"io"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/sync/errgroup"
)

type server struct {
	impl  ContainerRuntime
	stdio *plugin.StdioServer
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

func (s *server) GetStreamingConnection(_ context.Context, _ *api.GetStreamingConnectionRequest) (*api.GetStreamingConnectionResponse, error) {
	stdio, err := plugin.NewStdioServer()
	if err != nil {
		return nil, fmt.Errorf("could not get stdio server: %w", err)
	}

	s.stdio = stdio

	return &api.GetStreamingConnectionResponse{Url: stdio.Endpoint()}, nil
}

func (s *server) StreamContainer(_ context.Context, request *api.StreamContainerRequest) (*api.StreamContainerResponse, error) {
	inReader, inWriter := io.Pipe()
	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return s.stdio.Copy(inWriter, outReader, errReader)
	})

	eg.Go(func() error {
		defer func() {
			inWriter.Close()
			outWriter.Close()
			errWriter.Close()
		}()

		return s.impl.StreamContainer(request.ContainerId, &plugin.StreamConfig{
			Stdin:          inReader,
			Stdout:         outWriter,
			Stderr:         errWriter,
			TerminalHeight: request.Height,
			TerminalWidth:  request.Width,
		})
	})

	err := eg.Wait()

	result := &Result{}
	if ok := errors.As(err, result); ok {
		return &api.StreamContainerResponse{
			ExitCode: result.ExitCode,
			Message:  result.Message,
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error during container stream: %w", err)
	}

	return &api.StreamContainerResponse{ExitCode: 0}, nil
}
