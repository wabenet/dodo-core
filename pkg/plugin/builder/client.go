package builder

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/build/v1alpha2"
	pluginapi "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

const lenStreamID = 32

var _ ImageBuilder = &Client{}

type Client struct {
	builderClient api.PluginClient
	stdout        *grpcutil.StreamOutputClient
}

func NewGRPCClient(conn grpc.ClientConnInterface) *Client {
	return &Client{
		builderClient: api.NewPluginClient(conn),
		stdout:        grpcutil.NewStreamOutputClient(),
	}
}

func (c *Client) Type() plugin.Type { //nolint:ireturn
	return Type
}

func (c *Client) Metadata() plugin.Metadata {
	info, err := c.builderClient.GetPluginMetadata(context.Background(), &empty.Empty{})
	if err != nil {
		return plugin.NewFailedPluginInfo(Type, err)
	}

	return plugin.MetadataFromProto(info)
}

func (c *Client) Init() (plugin.Config, error) {
	resp, err := c.builderClient.InitPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return resp.GetConfig(), nil
}

func (c *Client) Cleanup() {
	_, err := c.builderClient.ResetPlugin(context.Background(), &empty.Empty{})
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

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error { return c.copyOutputClientToStdout(streamID, stream.Stdout, stream.Stderr) })

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

func (c *Client) copyOutputClientToStdout(streamID string, stdout, stderr io.Writer) error {
	req := &pluginapi.StreamOutputRequest{}

	req.SetId(streamID)

	outputClient, err := c.builderClient.StreamOutput(context.Background(), req)
	if err != nil {
		return fmt.Errorf("could not stream runtime output: %w", err)
	}

	if err := c.stdout.StreamOutput(outputClient, stdout, stderr); err != nil {
		return fmt.Errorf("could not stream runtime output: %w", err)
	}

	return nil
}
