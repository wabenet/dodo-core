package runtime

import (
	"context"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/sync/errgroup"
)

var _ ContainerRuntime = &client{}

type client struct {
	runtimeClient api.RuntimePluginClient
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) PluginInfo() (*api.PluginInfo, error) {
	return c.runtimeClient.GetPluginInfo(context.Background(), &empty.Empty{})
}

func (c *client) ResolveImage(spec string) (string, error) {
	img, err := c.runtimeClient.GetImage(context.Background(), &api.GetImageRequest{ImageSpec: spec})
	if err != nil {
		return "", err
	}

	return img.ImageId, nil
}

func (c *client) CreateContainer(config *api.Backdrop, tty bool, stdio bool) (string, error) {
	resp, err := c.runtimeClient.CreateContainer(context.Background(), &api.CreateContainerRequest{
		Config: config,
		Tty:    tty,
		Stdio:  stdio,
	})
	if err != nil {
		return "", err
	}

	return resp.ContainerId, nil
}

func (c *client) StartContainer(id string) error {
	_, err := c.runtimeClient.StartContainer(context.Background(), &api.StartContainerRequest{ContainerId: id})
	return err
}

func (c *client) DeleteContainer(id string) error {
	_, err := c.runtimeClient.DeleteContainer(context.Background(), &api.DeleteContainerRequest{ContainerId: id})
	return err
}

func (c *client) ResizeContainer(id string, height uint32, width uint32) error {
	_, err := c.runtimeClient.ResizeContainer(
		context.Background(),
		&api.ResizeContainerRequest{ContainerId: id, Height: height, Width: width},
	)

	return err
}

func (c *client) StreamContainer(id string, stream *plugin.StreamConfig) error {
	ctx := context.Background()

	connInfo, err := c.runtimeClient.GetStreamingConnection(ctx, &api.GetStreamingConnectionRequest{})
	if err != nil {
		return err
	}

	stdio, err := plugin.NewStdioClient(connInfo.Url)
	if err != nil {
		return err
	}

	eg, _ := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return stdio.Copy(stream.Stdin, stream.Stdout, stream.Stderr)
	})

	eg.Go(func() error {
		result, err := c.runtimeClient.StreamContainer(ctx, &api.StreamContainerRequest{ContainerId: id, Height: stream.TerminalHeight, Width: stream.TerminalWidth})
		if err != nil {
			return err
		}

		return &Result{
			ExitCode: result.ExitCode,
			Message:  result.Message,
		}
	})

	return eg.Wait()
}
