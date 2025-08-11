package stdio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	log "github.com/hashicorp/go-hclog"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OutputStreamingClient struct {
	streamOutputClient api.OutputStreamingPluginClient
	stdout             *StreamOutputClient
}

func NewOutputStreamingClient(conn grpc.ClientConnInterface) *OutputStreamingClient {
	return &OutputStreamingClient{
		streamOutputClient: api.NewOutputStreamingPluginClient(conn),
		stdout:             NewStreamOutputClient(),
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

type StreamOutputClient struct{}

type grpcOutputClient interface {
	Recv() (*api.StreamOutputResponse, error)
}

func NewStreamOutputClient() *StreamOutputClient {
	return &StreamOutputClient{}
}

func (*StreamOutputClient) StreamOutput(cl grpcOutputClient, stdout, stderr io.Writer) error {
	for {
		data, err := cl.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) ||
				errors.Is(err, context.Canceled) ||
				status.Code(err) == codes.Unavailable ||
				status.Code(err) == codes.Canceled ||
				status.Code(err) == codes.Unimplemented {
				return nil
			}

			return fmt.Errorf("error receiving data: %w", err)
		}

		switch data.GetChannel() {
		case api.OutputChannel_OUTPUT_CHANNEL_STDOUT:
			if _, err := io.Copy(stdout, bytes.NewReader(data.GetData())); err != nil {
				log.L().Error("failed to copy all bytes", "err", err)
			}

		case api.OutputChannel_OUTPUT_CHANNEL_STDERR:
			if _, err := io.Copy(stderr, bytes.NewReader(data.GetData())); err != nil {
				log.L().Error("failed to copy all bytes", "err", err)
			}

		case api.OutputChannel_OUTPUT_CHANNEL_UNSPECIFIED:
			log.L().Warn("unknown channel, dropping", "channel", data.GetChannel())

			continue
		}
	}
}
