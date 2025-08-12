package runtime

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	pluginapi "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/runtime/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/stdio"
	"golang.org/x/sync/errgroup"
)

var ErrUnexpectedMapType = errors.New("unexpected map type for stdio streaming server")

type Server struct {
	pluginapi.UnsafePluginServer
	pluginapi.UnsafeInputStreamingPluginServer
	api.UnsafeRuntimePluginServer

	stdio.OutputStreamingPluginServer
	stdio.InputStreamingPluginServer

	impl ContainerRuntime
}

func NewGRPCServer(impl ContainerRuntime) *Server {
	return &Server{impl: impl}
}

func (s *Server) reset() {
	s.OutputStreamingPluginServer.Reset()
	s.InputStreamingPluginServer.Reset()
}

func (s *Server) GetPluginMetadata(_ context.Context, _ *empty.Empty) (*pluginapi.GetPluginMetadataResponse, error) {
	resp := &pluginapi.GetPluginMetadataResponse{}

	resp.SetMetadata(s.impl.Metadata().ToProto())

	return resp, nil
}

func (s *Server) InitPlugin(_ context.Context, _ *empty.Empty) (*pluginapi.InitPluginResponse, error) {
	s.reset()

	config, err := s.impl.Init()
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	pluginConfig := &pluginapi.PluginConfig{}
	resp := &pluginapi.InitPluginResponse{}

	pluginConfig.SetConfig(config)
	resp.SetConfig(pluginConfig)

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

func (s *Server) StreamContainer(
	_ context.Context,
	request *api.StreamContainerRequest,
) (*api.StreamContainerResponse, error) {
	resp := &api.StreamContainerResponse{}

	inputStream, err := s.InputStreamingPluginServer.NewServerInputStream(request.GetContainerId())
	if err != nil {
		return nil, err
	}

	outputStream, err := s.OutputStreamingPluginServer.NewServerOutputStream(request.GetContainerId())
	if err != nil {
		return nil, err
	}

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(inputStream.Copy)
	eg.Go(outputStream.Copy)

	eg.Go(func() error {
		defer outputStream.Close()
		defer inputStream.Close()

		streamResp, err := s.impl.StreamContainer(request.GetContainerId(), &plugin.StreamConfig{
			Stdin:          inputStream.Stdin,
			Stdout:         outputStream.Stdout,
			Stderr:         outputStream.Stderr,
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
