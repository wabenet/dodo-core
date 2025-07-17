package plugin

import (
	"context"
	"fmt"
	"io"

	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
	"google.golang.org/grpc"
)

type OutputStreamingClient struct {
	streamOutputClient api.OutputStreamingPluginClient
	stdout             *grpcutil.StreamOutputClient
}

func NewOutputStreamingClient(conn grpc.ClientConnInterface) *OutputStreamingClient {
	return &OutputStreamingClient{
		streamOutputClient: api.NewOutputStreamingPluginClient(conn),
		stdout:             grpcutil.NewStreamOutputClient(),
	}
}

func (c *OutputStreamingClient) copyOutputClientToStdout(streamID string, stdout, stderr io.Writer) error {
	req := &api.StreamOutputRequest{}

	req.SetId(streamID)

	outputClient, err := c.streamOutputClient.StreamOutput(context.Background(), req)
	if err != nil {
		return fmt.Errorf("could not stream runtime output: %w", err)
	}

	if err := c.stdout.StreamOutput(outputClient, stdout, stderr); err != nil {
		return fmt.Errorf("could not stream runtime output: %w", err)
	}

	return nil
}

type ClientOutputStream struct {
	Copy func() error
}

func (c *OutputStreamingClient) PrepareStream(streamID string, stdout, stderr io.Writer) (ClientOutputStream, error) {
	return ClientOutputStream{
		Copy: func() error {
			return c.copyOutputClientToStdout(streamID, stdout, stderr)
		},
	}, nil
}
