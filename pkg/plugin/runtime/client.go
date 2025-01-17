package runtime

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	core "github.com/wabenet/dodo-core/api/core/v1alpha7"
	pluginapi "github.com/wabenet/dodo-core/api/plugin/v1alpha1"
	runtime "github.com/wabenet/dodo-core/api/runtime/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
	"github.com/wabenet/dodo-core/pkg/ioutil"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var _ ContainerRuntime = &client{}

type client struct {
	runtimeClient runtime.PluginClient
	stdin         *grpcutil.StreamInputClient
	stdout        *grpcutil.StreamOutputClient
}

func NewGRPCClient(conn grpc.ClientConnInterface) ContainerRuntime {
	return &client{
		runtimeClient: runtime.NewPluginClient(conn),
		stdin:         grpcutil.NewStreamInputClient(),
		stdout:        grpcutil.NewStreamOutputClient(),
	}
}

type streamInputClient struct {
	client runtime.Plugin_StreamInputClient
}

func (s *streamInputClient) Send(data *pluginapi.InputData) error {
	if err := s.client.Send(&pluginapi.StreamInputRequest{
		InputRequestType: &pluginapi.StreamInputRequest_InputData{InputData: data},
	}); err != nil {
		return fmt.Errorf("error wrapping Send call: %w", err)
	}

	return nil
}

func (s *streamInputClient) CloseAndRecv() (*empty.Empty, error) {
	e, err := s.client.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("error wrapping CloseAndRecv call: %w", err)
	}

	return e, nil
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) PluginInfo() *pluginapi.PluginInfo {
	info, err := c.runtimeClient.GetPluginInfo(context.Background(), &empty.Empty{})
	if err != nil {
		return plugin.NewFailedPluginInfo(Type, err)
	}

	return info
}

func (c *client) Init() (plugin.Config, error) {
	resp, err := c.runtimeClient.InitPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return resp.GetConfig(), nil
}

func (c *client) Cleanup() {
	_, err := c.runtimeClient.ResetPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		log.L().Error("plugin reset error", "error", err)
	}
}

func (c *client) ResolveImage(spec string) (string, error) {
	img, err := c.runtimeClient.GetImage(context.Background(), &runtime.GetImageRequest{ImageSpec: spec})
	if err != nil {
		return "", fmt.Errorf("could not resolve image: %w", err)
	}

	return img.GetImageId(), nil
}

func (c *client) CreateContainer(config *core.Backdrop, tty, stdio bool) (string, error) {
	resp, err := c.runtimeClient.CreateContainer(context.Background(), &runtime.CreateContainerRequest{
		Config: config,
		Tty:    tty,
		Stdio:  stdio,
	})
	if err != nil {
		return "", fmt.Errorf("could not create container: %w", err)
	}

	return resp.GetContainerId(), nil
}

func (c *client) StartContainer(id string) error {
	if _, err := c.runtimeClient.StartContainer(
		context.Background(),
		&runtime.StartContainerRequest{ContainerId: id},
	); err != nil {
		return fmt.Errorf("could not start container: %w", err)
	}

	return nil
}

func (c *client) DeleteContainer(id string) error {
	if _, err := c.runtimeClient.DeleteContainer(
		context.Background(),
		&runtime.DeleteContainerRequest{ContainerId: id},
	); err != nil {
		return fmt.Errorf("could not delete container: %w", err)
	}

	return nil
}

func (c *client) ResizeContainer(id string, height, width uint32) error {
	if _, err := c.runtimeClient.ResizeContainer(
		context.Background(),
		&runtime.ResizeContainerRequest{ContainerId: id, Height: height, Width: width},
	); err != nil {
		return fmt.Errorf("could not resize container: %w", err)
	}

	return nil
}

func (c *client) KillContainer(id string, signal os.Signal) error {
	if _, err := c.runtimeClient.KillContainer(
		context.Background(),
		&runtime.KillContainerRequest{ContainerId: id, Signal: signalToString(signal)},
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

	eg.Go(func() error { return c.copyInputClientToStdin(id, inReader) })
	eg.Go(func() error { return c.copyOutputClientToStdout(id, stream.Stdout, stream.Stderr) })

	eg.Go(func() error {
		defer inCancel()

		r, err := c.runtimeClient.StreamContainer(context.Background(), &runtime.StreamContainerRequest{
			ContainerId: id,
			Height:      stream.TerminalHeight,
			Width:       stream.TerminalWidth,
		})
		if err != nil {
			return fmt.Errorf("could not stream container: %w", err)
		}

		result.ExitCode = int(r.GetExitCode())

		return nil
	})

	if err := eg.Wait(); err != nil {
		return result, fmt.Errorf("error during container stream: %w", err)
	}

	return result, nil
}

func (c *client) copyInputClientToStdin(containerID string, stdin io.Reader) error {
	inputClient, err := c.runtimeClient.StreamInput(context.Background())
	if err != nil {
		return fmt.Errorf("could not stream runtime input: %w", err)
	}

	if err := inputClient.Send(&pluginapi.StreamInputRequest{
		InputRequestType: &pluginapi.StreamInputRequest_InitialRequest{
			InitialRequest: &pluginapi.InitialStreamInputRequest{Id: containerID},
		},
	}); err != nil {
		return fmt.Errorf("could not stream runtime input: %w", err)
	}

	if err := c.stdin.StreamInput(&streamInputClient{client: inputClient}, stdin); err != nil {
		return fmt.Errorf("could not stream runtime input: %w", err)
	}

	return nil
}

func (c *client) copyOutputClientToStdout(containerID string, stdout, stderr io.Writer) error {
	outputClient, err := c.runtimeClient.StreamOutput(
		context.Background(),
		&pluginapi.StreamOutputRequest{Id: containerID},
	)
	if err != nil {
		return fmt.Errorf("could not stream runtime output: %w", err)
	}

	if err := c.stdout.StreamOutput(outputClient, stdout, stderr); err != nil {
		return fmt.Errorf("could not stream runtime output: %w", err)
	}

	return nil
}

func (c *client) CreateVolume(name string) error {
	if _, err := c.runtimeClient.CreateVolume(
		context.Background(),
		&runtime.CreateVolumeRequest{Name: name},
	); err != nil {
		return fmt.Errorf("could not create volume: %w", err)
	}

	return nil
}

func (c *client) DeleteVolume(name string) error {
	if _, err := c.runtimeClient.DeleteVolume(
		context.Background(),
		&runtime.DeleteVolumeRequest{Name: name},
	); err != nil {
		return fmt.Errorf("could not delete volume: %w", err)
	}

	return nil
}

func (c *client) WriteFile(containerID, path string, contents []byte) error {
	if _, err := c.runtimeClient.WriteFile(
		context.Background(),
		&runtime.WriteFileRequest{
			ContainerId: containerID,
			FilePath:    path,
			Contents:    string(contents),
		},
	); err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}
