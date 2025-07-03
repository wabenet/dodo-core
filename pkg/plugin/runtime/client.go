package runtime

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	pluginapi "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/runtime/v1alpha2"
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
	req := &pluginapi.StreamInputRequest{}

	req.SetInputData(data)

	if err := s.client.Send(req); err != nil {
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

func (c *Client) Metadata() plugin.Metadata {
	metadata, err := c.runtimeClient.GetPluginMetadata(context.Background(), &empty.Empty{})
	if err != nil {
		return plugin.NewFailedPluginInfo(Type, err)
	}

	return plugin.MetadataFromProto(metadata)
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
	req := &api.GetImageRequest{}

	req.SetImageSpec(spec)

	img, err := c.runtimeClient.GetImage(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("could not resolve image: %w", err)
	}

	return img.GetImageId(), nil
}

func (c *Client) CreateContainer(config ContainerConfig) (string, error) {
	req := &api.CreateContainerRequest{}

	req.SetConfig(config.ToProto())

	resp, err := c.runtimeClient.CreateContainer(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("could not create container: %w", err)
	}

	return resp.GetContainerId(), nil
}

func (c *Client) StartContainer(id string) error {
	req := &api.StartContainerRequest{}

	req.SetContainerId(id)

	if _, err := c.runtimeClient.StartContainer(context.Background(), req); err != nil {
		return fmt.Errorf("could not start container: %w", err)
	}

	return nil
}

func (c *Client) DeleteContainer(id string) error {
	req := &api.DeleteContainerRequest{}

	req.SetContainerId(id)

	if _, err := c.runtimeClient.DeleteContainer(context.Background(), req); err != nil {
		return fmt.Errorf("could not delete container: %w", err)
	}

	return nil
}

func (c *Client) ResizeContainer(id string, height, width uint32) error {
	req := &api.ResizeContainerRequest{}

	req.SetContainerId(id)
	req.SetHeight(height)
	req.SetWidth(width)

	if _, err := c.runtimeClient.ResizeContainer(context.Background(), req); err != nil {
		return fmt.Errorf("could not resize container: %w", err)
	}

	return nil
}

func (c *Client) KillContainer(id string, signal os.Signal) error {
	req := &api.KillContainerRequest{}

	req.SetContainerId(id)
	req.SetSignal(signalToString(signal))

	if _, err := c.runtimeClient.KillContainer(context.Background(), req); err != nil {
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

		req := &api.StreamContainerRequest{}

		req.SetContainerId(id)
		req.SetHeight(stream.TerminalHeight)
		req.SetWidth(stream.TerminalWidth)

		resp, err := c.runtimeClient.StreamContainer(context.Background(), req)
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

	req := &pluginapi.StreamInputRequest{}
	ir := &pluginapi.InitialStreamInputRequest{}

	ir.SetId(containerID)
	req.SetInitialRequest(ir)

	if err := inputClient.Send(req); err != nil {
		return fmt.Errorf("could not stream runtime input: %w", err)
	}

	if err := c.stdin.StreamInput(&streamInputClient{client: inputClient}, stdin); err != nil {
		return fmt.Errorf("could not stream runtime input: %w", err)
	}

	return nil
}

func (c *Client) copyOutputClientToStdout(containerID string, stdout, stderr io.Writer) error {
	req := &pluginapi.StreamOutputRequest{}

	req.SetId(containerID)

	outputClient, err := c.runtimeClient.StreamOutput(context.Background(), req)
	if err != nil {
		return fmt.Errorf("could not stream runtime output: %w", err)
	}

	if err := c.stdout.StreamOutput(outputClient, stdout, stderr); err != nil {
		return fmt.Errorf("could not stream runtime output: %w", err)
	}

	return nil
}

func (c *Client) CreateVolume(name string) error {
	req := &api.CreateVolumeRequest{}

	req.SetName(name)

	if _, err := c.runtimeClient.CreateVolume(context.Background(), req); err != nil {
		return fmt.Errorf("could not create volume: %w", err)
	}

	return nil
}

func (c *Client) DeleteVolume(name string) error {
	req := &api.DeleteVolumeRequest{}

	req.SetName(name)

	if _, err := c.runtimeClient.DeleteVolume(context.Background(), req); err != nil {
		return fmt.Errorf("could not delete volume: %w", err)
	}

	return nil
}

func (c *Client) WriteFile(containerID, path string, contents []byte) error {
	req := &api.WriteFileRequest{}

	req.SetContainerId(containerID)
	req.SetFilePath(path)
	req.SetContents(string(contents))

	if _, err := c.runtimeClient.WriteFile(context.Background(), req); err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}
