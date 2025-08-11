package runtime

import (
	"context"
	"fmt"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	pluginapi "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/runtime/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/stdio"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

var _ ContainerRuntime = &Client{}

type Client struct {
	pluginClient  pluginapi.PluginClient
	runtimeClient api.RuntimePluginClient

	streamOutputClient *stdio.OutputStreamingClient
	streamInputClient  *stdio.InputStreamingClient
}

func NewGRPCClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		pluginClient:  pluginapi.NewPluginClient(conn),
		runtimeClient: api.NewRuntimePluginClient(conn),

		streamOutputClient: stdio.NewOutputStreamingClient(conn),
		streamInputClient:  stdio.NewInputStreamingClient(conn),
	}
}

func (c *Client) Type() plugin.Type { //nolint:ireturn
	return Type
}

func (c *Client) Metadata() plugin.Metadata {
	resp, err := c.pluginClient.GetPluginMetadata(context.Background(), &empty.Empty{})
	if err != nil {
		return plugin.NewFailedPluginInfo(Type, err)
	}

	return plugin.MetadataFromProto(resp.GetMetadata())
}

func (c *Client) Init() (plugin.Config, error) {
	resp, err := c.pluginClient.InitPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return resp.GetConfig().GetConfig(), nil
}

func (c *Client) Cleanup() {
	_, err := c.pluginClient.ResetPlugin(context.Background(), &empty.Empty{})
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

	outputStream, err := c.streamOutputClient.PrepareStream(id, stream.Stdout, stream.Stderr)
	if err != nil {
		return nil, err
	}

	inputStream, err := c.streamInputClient.PrepareStream(id, stream.Stdin)
	if err != nil {
		return nil, err
	}

	eg.Go(inputStream.Copy)
	eg.Go(outputStream.Copy)

	eg.Go(func() error {
		defer inputStream.Close()

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
