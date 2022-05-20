package runtime

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	api "github.com/dodo-cli/dodo-core/api/v1alpha3"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
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

func (c *client) CreateContainer(config *api.Backdrop, tty bool, stdio bool) (string, error) {
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

func (c *client) ResizeContainer(id string, height uint32, width uint32) error {
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // TODO: replace the cancel with Close and Send/Recv pair on the stream thing?

	inputClient, err := c.runtimeClient.StreamRuntimeInput(ctx)
	if err != nil {
		return nil, err
	}

	outputClient, err := c.runtimeClient.StreamRuntimeOutput(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	go streamInput(inputClient, stream.Stdin)
	go streamOutput(outputClient, stream.Stdout, stream.Stderr)

	result, err := c.runtimeClient.StreamContainer(ctx, &api.StreamContainerRequest{
		ContainerId: id,
		Height:      stream.TerminalHeight,
		Width:       stream.TerminalWidth,
	})
	if err != nil {
		return nil, fmt.Errorf("could not stream container: %w", err)
	}

	return &Result{ExitCode: int(result.ExitCode)}, nil
}

func streamInput(c api.RuntimePlugin_StreamRuntimeInputClient, stdin io.Reader) {
	bufsrc := bufio.NewReader(stdin)
	data := api.InputData{}

	for {
		var b [1024]byte

		n, err := bufsrc.Read(b[:])

		if n > 0 {
			data.Data = b[:n]
			if serr := c.Send(&data); err != nil {
				log.L().Warn("error in stdio stream", "err", serr)
				return
			}
		}

		if err == io.EOF {
			c.CloseAndRecv()
			return
		}

		if err != nil {
			log.L().Warn("error in stdio stream", "err", err)
			return
		}
	}
}

func streamOutput(c api.RuntimePlugin_StreamRuntimeOutputClient, stdout io.Writer, stderr io.Writer) {
	for {
		data, err := c.Recv()
		if err != nil {
			if err == io.EOF ||
				status.Code(err) == codes.Unavailable ||
				status.Code(err) == codes.Canceled ||
				status.Code(err) == codes.Unimplemented ||
				err == context.Canceled {
				return
			}

			log.L().Error("error receiving data", "err", err)
			return
		}

		var w io.Writer
		switch data.Channel {
		case api.OutputData_STDOUT:
			w = stdout

		case api.OutputData_STDERR:
			w = stderr

		default:
			log.L().Warn("unknown channel, dropping", "channel", data.Channel)
			continue
		}

		if _, err := io.Copy(w, bytes.NewReader(data.Data)); err != nil {
			log.L().Error("failed to copy all bytes", "err", err)
		}
	}
}
