package builder

import (
	"context"
	"io"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/sync/errgroup"
)

type server struct {
	impl  ImageBuilder
	stdio *plugin.StdioServer
}

func (s *server) GetPluginInfo(_ context.Context, _ *empty.Empty) (*api.PluginInfo, error) {
	return s.impl.PluginInfo()
}

func (s *server) GetStreamingConnection(_ context.Context, _ *api.GetStreamingConnectionRequest) (*api.GetStreamingConnectionResponse, error) {
	stdio, err := plugin.NewStdioServer()
	if err != nil {
		return nil, err
	}

	s.stdio = stdio

	return &api.GetStreamingConnectionResponse{Url: stdio.Endpoint()}, nil
}

func (s *server) CreateImage(_ context.Context, request *api.CreateImageRequest) (*api.CreateImageResponse, error) {
	resp := &api.CreateImageResponse{}

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

		imageID, err := s.impl.CreateImage(request.Config, &plugin.StreamConfig{
			Stdin:          inReader,
			Stdout:         outWriter,
			Stderr:         errWriter,
			TerminalHeight: request.Height,
			TerminalWidth:  request.Width,
		})

		if err != nil {
			return err
		}

		resp.ImageId = imageID
		return nil
	})

	err := eg.Wait()

	if err != nil {
		return nil, err
	}

	return resp, nil
}
