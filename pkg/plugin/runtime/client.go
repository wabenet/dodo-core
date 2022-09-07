package runtime

import (
	"context"
	"fmt"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	api "github.com/wabenet/dodo-core/api/v1alpha4"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
	"github.com/wabenet/dodo-core/pkg/ioutil"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
)

var _ ContainerRuntime = &client{}

type client struct {
	runtimeClient api.RuntimePluginClient
	stdin         *grpcutil.StreamInputClient
	stdout        *grpcutil.StreamOutputClient
}

func NewGRPCClient(c api.RuntimePluginClient) ContainerRuntime {
	return &client{
		runtimeClient: c,
		stdin:         grpcutil.NewStreamInputClient(),
		stdout:        grpcutil.NewStreamOutputClient(),
	}
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) PluginInfo() *api.PluginInfo {
	info, err := c.runtimeClient.GetPluginInfo(context.Background(), &empty.Empty{})
	if err != nil {
		return &api.PluginInfo{
			Name:   &api.PluginName{Type: Type.String(), Name: plugin.FailedPlugin},
			Fields: map[string]string{"error": err.Error()},
		}
	}

	return info
}

func (c *client) Init() (plugin.PluginConfig, error) {
	resp, err := c.runtimeClient.InitPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return resp.Config, nil
}

func (c *client) Cleanup() {
	_, err := c.runtimeClient.ResetPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		log.L().Error("plugin reset error", "error", err)
	}
}

func (c *client) ResolveImage(spec string) (string, error) {
	img, err := c.runtimeClient.GetImage(context.Background(), &api.GetImageRequest{ImageSpec: spec})
	if err != nil {
		return "", fmt.Errorf("could not resolve image: %w", err)
	}

	return img.ImageId, nil
}

func (c *client) CreateContainer(config *api.Backdrop, tty, stdio bool) (string, error) {
	resp, err := c.runtimeClient.CreateContainer(context.Background(), &api.CreateContainerRequest{
		Config: config,
		Tty:    tty,
		Stdio:  stdio,
	})
	if err != nil {
		return "", fmt.Errorf("could not create container: %w", err)
	}

	return resp.ContainerId, nil
}

func (c *client) StartContainer(id string) error {
	if _, err := c.runtimeClient.StartContainer(
		context.Background(),
		&api.StartContainerRequest{ContainerId: id},
	); err != nil {
		return fmt.Errorf("could not start container: %w", err)
	}

	return nil
}

func (c *client) DeleteContainer(id string) error {
	if _, err := c.runtimeClient.DeleteContainer(
		context.Background(),
		&api.DeleteContainerRequest{ContainerId: id},
	); err != nil {
		return fmt.Errorf("could not delete container: %w", err)
	}

	return nil
}

func (c *client) ResizeContainer(id string, height, width uint32) error {
	if _, err := c.runtimeClient.ResizeContainer(
		context.Background(),
		&api.ResizeContainerRequest{ContainerId: id, Height: height, Width: width},
	); err != nil {
		return fmt.Errorf("could not resize container: %w", err)
	}

	return nil
}

func (c *client) KillContainer(id string, signal os.Signal) error {
	if _, err := c.runtimeClient.KillContainer(
		context.Background(),
		&api.KillContainerRequest{ContainerId: id, Signal: signalToString(signal)},
	); err != nil {
		return fmt.Errorf("could not kill container: %w", err)
	}

	return nil
}

func (c *client) StreamContainer(id string, stream *plugin.StreamConfig) (*Result, error) {
	result := &Result{}
	eg, _ := errgroup.WithContext(context.Background())

	inContext, inCancel := context.WithCancel(context.Background())
	inReader := ioutil.NewCancelableReader(inContext, stream.Stdin)

	eg.Go(func() error {
		inputClient, err := c.runtimeClient.StreamInput(context.Background())
		if err != nil {
			return fmt.Errorf("could not stream runtime input: %w", err)
		}

		if err := c.stdin.StreamInput(inputClient, inReader); err != nil {
			return fmt.Errorf("could not stream runtime input: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		outputClient, err := c.runtimeClient.StreamOutput(context.Background(), &empty.Empty{})
		if err != nil {
			return fmt.Errorf("could not stream runtime output: %w", err)
		}

		if err := c.stdout.StreamOutput(outputClient, stream.Stdout, stream.Stderr); err != nil {
			return fmt.Errorf("could not stream runtime output: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		defer inCancel()

		r, err := c.runtimeClient.StreamContainer(context.Background(), &api.StreamContainerRequest{
			ContainerId: id,
			Height:      stream.TerminalHeight,
			Width:       stream.TerminalWidth,
		})
		if err != nil {
			return fmt.Errorf("could not stream container: %w", err)
		}

		result.ExitCode = int(r.ExitCode)

		return nil
	})

	if err := eg.Wait(); err != nil {
		return result, fmt.Errorf("error during container stream: %w", err)
	}

	return result, nil
}
