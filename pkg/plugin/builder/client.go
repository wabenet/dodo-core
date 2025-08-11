package builder

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/build/v1alpha2"
	pluginapi "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/stdio"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

const lenStreamID = 32

var _ ImageBuilder = &Client{}

type Client struct {
	pluginClient  pluginapi.PluginClient
	builderClient api.BuilderPluginClient

	streamOutputClient *stdio.OutputStreamingClient
}

func NewGRPCClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		pluginClient:  pluginapi.NewPluginClient(conn),
		builderClient: api.NewBuilderPluginClient(conn),

		streamOutputClient: stdio.NewOutputStreamingClient(conn),
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

func (c *Client) CreateImage(config BuildConfig, stream *plugin.StreamConfig) (string, error) {
	if stream == nil {
		req := &api.CreateImageRequest{}

		req.SetConfig(config.ToProto())

		result, err := c.builderClient.CreateImage(context.Background(), req)
		if err != nil {
			return "", fmt.Errorf("could not build image: %w", err)
		}

		return result.GetImageId(), nil
	}

	b := make([]byte, lenStreamID)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("could not generate stream id: %w", err)
	}

	streamID := hex.EncodeToString(b)
	imageID := ""
	req := &api.CreateImageRequest{}

	req.SetStreamId(streamID)
	req.SetConfig(config.ToProto())
	req.SetHeight(stream.TerminalHeight)
	req.SetWidth(stream.TerminalWidth)

	outputStream, err := c.streamOutputClient.PrepareStream(streamID, stream.Stdout, stream.Stderr)
	if err != nil {
		return "", err
	}

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(outputStream.Copy)

	eg.Go(func() error {
		result, err := c.builderClient.CreateImage(context.Background(), req)
		if err != nil {
			return fmt.Errorf("could not build image: %w", err)
		}

		imageID = result.GetImageId()

		return nil
	})

	if err := eg.Wait(); err != nil {
		return "", fmt.Errorf("error during image build stream: %w", err)
	}

	return imageID, nil
}
