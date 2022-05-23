package builder

import (
	"bytes"
	"context"
	"fmt"
	"io"

	api "github.com/dodo-cli/dodo-core/api/v1alpha3"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ ImageBuilder = &client{}

type client struct {
	builderClient api.BuilderPluginClient
}

func NewGRPCClient(c api.BuilderPluginClient) ImageBuilder {
	return &client{builderClient: c}
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) PluginInfo() *api.PluginInfo {
	info, err := c.builderClient.GetPluginInfo(context.Background(), &empty.Empty{})
	if err != nil {
		return &api.PluginInfo{
			Name:   &api.PluginName{Type: Type.String(), Name: plugin.FailedPlugin},
			Fields: map[string]string{"error": err.Error()},
		}
	}

	return info
}

func (c *client) Init() (plugin.PluginConfig, error) {
	resp, err := c.builderClient.InitPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return resp.Config, nil
}

func (c *client) CreateImage(config *api.BuildInfo, stream *plugin.StreamConfig) (string, error) {
	if stream == nil {
		result, err := c.builderClient.CreateImage(context.Background(), &api.CreateImageRequest{
			Config: config,
			Height: 0,
			Width:  0,
		})

		if err != nil {
			return "", fmt.Errorf("could not build image: %w", err)
		}

		return result.ImageId, nil
	}

	stdioClient, err := c.builderClient.StreamBuildOutput(context.Background(), &empty.Empty{})
	if err != nil {
		return "", fmt.Errorf("could not stream build output: %w", err)
	}

	go streamOutput(stdioClient, stream.Stdout, stream.Stderr)

	result, err := c.builderClient.CreateImage(context.Background(), &api.CreateImageRequest{
		Config: config,
		Height: stream.TerminalHeight,
		Width:  stream.TerminalWidth,
	})
	if err != nil {
		return "", fmt.Errorf("could not build image: %w", err)
	}

	return result.ImageId, nil
}

func streamOutput(c api.BuilderPlugin_StreamBuildOutputClient, stdout io.Writer, stderr io.Writer) {
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
