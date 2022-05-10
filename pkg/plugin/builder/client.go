package builder

import (
	"context"
	"fmt"

	api "github.com/dodo-cli/dodo-core/api/v1alpha3"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/sync/errgroup"
)

var _ ImageBuilder = &client{}

type client struct {
	builderClient api.BuilderPluginClient
}

func NewGRPCClient(c api.BuilderPluginClient) ImageBuilder {
	return &client{builderClient: c}
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) PluginInfo() *api.PluginInfo {
	info, err := c.builderClient.GetPluginInfo(context.Background(), &empty.Empty{})
	if err != nil {
		return &api.PluginInfo{
			Name:   &api.PluginName{Type: Type.String(), Name: plugin.FailedPlugin},
			Fields: map[string]string{"error": err.Error()},
		}
	}

	return info
}

func (c *client) Init() (plugin.PluginConfig, error) {
	resp, err := c.builderClient.InitPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return resp.Config, nil
}

func (c *client) CreateImage(config *api.BuildInfo, stream *plugin.StreamConfig) (string, error) {
	ctx := context.Background()
	imageID := ""

	connInfo, err := c.builderClient.GetStreamingConnection(ctx, &api.GetStreamingConnectionRequest{})
	if err != nil {
		return "", fmt.Errorf("could not get streaming connection: %w", err)
	}

	stdio, err := plugin.NewStdioClient(connInfo.Url)
	if err != nil {
		return "", fmt.Errorf("could not get stdio client: %w", err)
	}

	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return stdio.Copy(stream.Stdin, stream.Stdout, stream.Stderr)
	})

	eg.Go(func() error {
		result, err := c.builderClient.CreateImage(ctx, &api.CreateImageRequest{Config: config, Height: stream.TerminalHeight, Width: stream.TerminalWidth})
		if err != nil {
			return fmt.Errorf("could not build image: %w", err)
		}

		imageID = result.ImageId

		return nil
	})

	return imageID, eg.Wait()
}
