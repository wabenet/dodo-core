package runtime

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	api "github.com/wabenet/dodo-core/api/v1alpha3"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ ContainerRuntime = &client{}

type client struct {
	runtimeClient api.RuntimePluginClient
}

func NewGRPCClient(c api.RuntimePluginClient) ContainerRuntime {
	return &client{runtimeClient: c}
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

	eg.Go(func() error {
		inputClient, err := c.runtimeClient.StreamRuntimeInput(context.Background())
		if err != nil {
			return fmt.Errorf("could not stream runtime input: %w", err)
		}

		return streamInput(inputClient, stream.Stdin)
	})

	eg.Go(func() error {
		outputClient, err := c.runtimeClient.StreamRuntimeOutput(context.Background(), &empty.Empty{})
		if err != nil {
			return fmt.Errorf("could not stream runtime output: %w", err)
		}

		return streamOutput(outputClient, stream.Stdout, stream.Stderr)
	})

	eg.Go(func() error {
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

func streamInput(c api.RuntimePlugin_StreamRuntimeInputClient, stdin io.Reader) error {
	bufsrc := bufio.NewReader(stdin)
	data := api.InputData{}

	for {
		var b [1024]byte

		n, err := bufsrc.Read(b[:])

		if n > 0 {
			data.Data = b[:n]
			if serr := c.Send(&data); err != nil {
				return fmt.Errorf("could not send input to server: %w", serr)
			}
		}

		if err == io.EOF {
			if _, serr := c.CloseAndRecv(); serr != nil {
				log.L().Warn("could not close input stream", "err", serr)
			}

			return nil
		}

		if err != nil {
			return fmt.Errorf("could not read stream from client: %w", err)
		}
	}
}

func streamOutput(c api.RuntimePlugin_StreamRuntimeOutputClient, stdout, stderr io.Writer) error {
	for {
		data, err := c.Recv()
		if err != nil {
			if err == io.EOF ||
				status.Code(err) == codes.Unavailable ||
				status.Code(err) == codes.Canceled ||
				status.Code(err) == codes.Unimplemented ||
				err == context.Canceled {
				return nil
			}

			return fmt.Errorf("error receiving data: %w", err)
		}

		switch data.Channel {
		case api.OutputData_STDOUT:
			if _, err := io.Copy(stdout, bytes.NewReader(data.Data)); err != nil {
				log.L().Error("failed to copy all bytes", "err", err)
			}

		case api.OutputData_STDERR:
			if _, err := io.Copy(stderr, bytes.NewReader(data.Data)); err != nil {
				log.L().Error("failed to copy all bytes", "err", err)
			}

		case api.OutputData_INVALID:
			log.L().Warn("unknown channel, dropping", "channel", data.Channel)

			continue
		}
	}
}
