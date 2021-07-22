package builder

import (
	"context"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/sync/errgroup"
)

var _ ImageBuilder = &client{}

type client struct {
	builderClient api.BuilderPluginClient
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) PluginInfo() (*api.PluginInfo, error) {
	return c.builderClient.GetPluginInfo(context.Background(), &empty.Empty{})
}

func (c *client) CreateImage(config *api.BuildInfo, stream *plugin.StreamConfig) (string, error) {
	ctx := context.Background()
	imageId := ""

	connInfo, err := c.builderClient.GetStreamingConnection(ctx, &api.GetStreamingConnectionRequest{})
	if err != nil {
		return "", err
	}

	stdio, err := plugin.NewStdioClient(connInfo.Url)
	if err != nil {
		return "", err
	}

	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return stdio.Copy(stream.Stdin, stream.Stdout, stream.Stderr)
	})

	eg.Go(func() error {
		result, err := c.builderClient.CreateImage(ctx, &api.CreateImageRequest{Config: config, Height: stream.TerminalHeight, Width: stream.TerminalWidth})
		if err != nil {
			return err
		}

		imageId = result.ImageId
		return nil
	})

	return imageId, eg.Wait()
}
