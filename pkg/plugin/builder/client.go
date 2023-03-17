package builder

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	build "github.com/wabenet/dodo-core/api/build/v1alpha1"
	core "github.com/wabenet/dodo-core/api/core/v1alpha5"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"golang.org/x/sync/errgroup"
)

const lenStreamID = 32

var _ ImageBuilder = &client{}

type client struct {
	builderClient build.PluginClient
	stdout        *grpcutil.StreamOutputClient
}

func NewGRPCClient(c build.PluginClient) ImageBuilder {
	return &client{
		builderClient: c,
		stdout:        grpcutil.NewStreamOutputClient(),
	}
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) PluginInfo() *core.PluginInfo {
	info, err := c.builderClient.GetPluginInfo(context.Background(), &empty.Empty{})
	if err != nil {
		return &core.PluginInfo{
			Name:   &core.PluginName{Type: Type.String(), Name: plugin.FailedPlugin},
			Fields: map[string]string{"error": err.Error()},
		}
	}

	return info
}

func (c *client) Init() (plugin.Config, error) {
	resp, err := c.builderClient.InitPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return resp.Config, nil
}

func (c *client) Cleanup() {
	_, err := c.builderClient.ResetPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		log.L().Error("plugin reset error", "error", err)
	}
}

func (c *client) CreateImage(config *core.BuildInfo, stream *plugin.StreamConfig) (string, error) {
	if stream == nil {
		result, err := c.builderClient.CreateImage(context.Background(), &build.CreateImageRequest{
			Config: config,
			Height: 0,
			Width:  0,
		})
		if err != nil {
			return "", fmt.Errorf("could not build image: %w", err)
		}

		return result.ImageId, nil
	}

	b := make([]byte, lenStreamID)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("could not generate stream id: %w", err)
	}

	streamID := hex.EncodeToString(b)
	imageID := ""

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error { return c.copyOutputClientToStdout(streamID, stream.Stdout, stream.Stderr) })

	eg.Go(func() error {
		result, err := c.builderClient.CreateImage(context.Background(), &build.CreateImageRequest{
			StreamId: streamID,
			Config:   config,
			Height:   stream.TerminalHeight,
			Width:    stream.TerminalWidth,
		})
		if err != nil {
			return fmt.Errorf("could not build image: %w", err)
		}

		imageID = result.ImageId

		return nil
	})

	if err := eg.Wait(); err != nil {
		return "", fmt.Errorf("error during image build stream: %w", err)
	}

	return imageID, nil
}

func (c *client) copyOutputClientToStdout(streamID string, stdout, stderr io.Writer) error {
	outputClient, err := c.builderClient.StreamOutput(
		context.Background(),
		&core.StreamOutputRequest{Id: streamID},
	)
	if err != nil {
		return fmt.Errorf("could not stream runtime output: %w", err)
	}

	if err := c.stdout.StreamOutput(outputClient, stdout, stderr); err != nil {
		return fmt.Errorf("could not stream runtime output: %w", err)
	}

	return nil
}
