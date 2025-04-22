package runtime

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	pluginapi "github.com/wabenet/dodo-core/api/plugin/v1alpha1"
	api "github.com/wabenet/dodo-core/api/runtime/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
	"github.com/wabenet/dodo-core/pkg/ioutil"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var _ ContainerRuntime = &Client{}

type Client struct {
	runtimeClient api.PluginClient
	stdin         *grpcutil.StreamInputClient
	stdout        *grpcutil.StreamOutputClient
}

func NewGRPCClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		runtimeClient: api.NewPluginClient(conn),
		stdin:         grpcutil.NewStreamInputClient(),
		stdout:        grpcutil.NewStreamOutputClient(),
	}
}

type streamInputClient struct {
	client api.Plugin_StreamInputClient
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

func (c *Client) Type() plugin.Type { //nolint:ireturn
	return Type
}

func (c *Client) PluginInfo() *pluginapi.PluginInfo {
	info, err := c.runtimeClient.GetPluginInfo(context.Background(), &empty.Empty{})
	if err != nil {
		return plugin.NewFailedPluginInfo(Type, err)
	}

	return info
}

func (c *Client) Init() (plugin.Config, error) {
	resp, err := c.runtimeClient.InitPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return resp.GetConfig(), nil
}

func (c *Client) Cleanup() {
	_, err := c.runtimeClient.ResetPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		log.L().Error("plugin reset error", "error", err)
	}
}

func (c *Client) ResolveImage(spec string) (string, error) {
	img, err := c.runtimeClient.GetImage(context.Background(), &api.GetImageRequest{ImageSpec: spec})
	if err != nil {
		return "", fmt.Errorf("could not resolve image: %w", err)
	}

	return img.GetImageId(), nil
}

func (c *Client) CreateContainer(config *api.ContainerConfig) (string, error) {
	resp, err := c.runtimeClient.CreateContainer(context.Background(), &api.CreateContainerRequest{Config: config})
	if err != nil {
		return "", fmt.Errorf("could not create container: %w", err)
	}

	return resp.GetContainerId(), nil
}

func (c *Client) StartContainer(id string) error {
	if _, err := c.runtimeClient.StartContainer(
		context.Background(),
		&api.StartContainerRequest{ContainerId: id},
	); err != nil {
		return fmt.Errorf("could not start container: %w", err)
	}

	return nil
}

func (c *Client) DeleteContainer(id string) error {
	if _, err := c.runtimeClient.DeleteContainer(
		context.Background(),
		&api.DeleteContainerRequest{ContainerId: id},
	); err != nil {
		return fmt.Errorf("could not delete container: %w", err)
	}

	return nil
}

func (c *Client) ResizeContainer(id string, height, width uint32) error {
	if _, err := c.runtimeClient.ResizeContainer(
		context.Background(),
		&api.ResizeContainerRequest{ContainerId: id, Height: height, Width: width},
	); err != nil {
		return fmt.Errorf("could not resize container: %w", err)
	}

	return nil
}

func (c *Client) KillContainer(id string, signal os.Signal) error {
	if _, err := c.runtimeClient.KillContainer(
		context.Background(),
		&api.KillContainerRequest{ContainerId: id, Signal: signalToString(signal)},
	); err != nil {
		return fmt.Errorf("could not kill container: %w", err)
	}

	return nil
}

func (c *Client) StreamContainer(id string, stream *plugin.StreamConfig) (*Result, error) {
	result := &Result{}
	eg, _ := errgroup.WithContext(context.Background())

	inContext, inCancel := context.WithCancel(context.Background())
	inReader := ioutil.NewCancelableReader(inContext, stream.Stdin)

	eg.Go(func() error { return c.copyInputClientToStdin(id, inReader) })
	eg.Go(func() error { return c.copyOutputClientToStdout(id, stream.Stdout, stream.Stderr) })

	eg.Go(func() error {
		defer inCancel()

		resp, err := c.runtimeClient.StreamContainer(context.Background(), &api.StreamContainerRequest{
			ContainerId: id,
			Height:      stream.TerminalHeight,
			Width:       stream.TerminalWidth,
		})
		if err != nil {
			return fmt.Errorf("could not stream container: %w", err)
		}

		result.ExitCode = int(resp.GetExitCode())

		return nil
	})

	if err := eg.Wait(); err != nil {
		return result, fmt.Errorf("error during container stream: %w", err)
	}

	return result, nil
}

func (c *Client) copyInputClientToStdin(containerID string, stdin io.Reader) error {
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

func (c *Client) copyOutputClientToStdout(containerID string, stdout, stderr io.Writer) error {
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

func (c *Client) CreateVolume(name string) error {
	if _, err := c.runtimeClient.CreateVolume(
		context.Background(),
		&api.CreateVolumeRequest{Name: name},
	); err != nil {
		return fmt.Errorf("could not create volume: %w", err)
	}

	return nil
}

func (c *Client) DeleteVolume(name string) error {
	if _, err := c.runtimeClient.DeleteVolume(
		context.Background(),
		&api.DeleteVolumeRequest{Name: name},
	); err != nil {
		return fmt.Errorf("could not delete volume: %w", err)
	}

	return nil
}

func (c *Client) WriteFile(containerID, path string, contents []byte) error {
	if _, err := c.runtimeClient.WriteFile(
		context.Background(),
		&api.WriteFileRequest{
			ContainerId: containerID,
			FilePath:    path,
			Contents:    string(contents),
		},
	); err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}
