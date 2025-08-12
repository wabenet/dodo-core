package stdio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	log "github.com/hashicorp/go-hclog"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClientOutputStream struct {
	client   api.OutputStreamingPluginClient
	streamID string
	stdout   io.Writer
	stderr   io.Writer
}

func NewClientOutputStream(client api.OutputStreamingPluginClient, streamID string, stdout, stderr io.Writer) (ClientOutputStream, error) {
	return ClientOutputStream{
		client:   client,
		streamID: streamID,
		stdout:   stdout,
		stderr:   stderr,
	}, nil
}

func (s ClientOutputStream) Copy() error {
	req := &api.StreamOutputRequest{}

	req.SetId(s.streamID)

	client, err := s.client.StreamOutput(context.Background(), req)
	if err != nil {
		return fmt.Errorf("could not stream runtime output: %w", err)
	}

	for {
		data, err := client.Recv()
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
			if _, err := io.Copy(s.stdout, bytes.NewReader(data.GetData())); err != nil {
				log.L().Error("failed to copy all bytes", "err", err)
			}

		case api.OutputChannel_OUTPUT_CHANNEL_STDERR:
			if _, err := io.Copy(s.stderr, bytes.NewReader(data.GetData())); err != nil {
				log.L().Error("failed to copy all bytes", "err", err)
			}

		case api.OutputChannel_OUTPUT_CHANNEL_UNSPECIFIED:
			log.L().Warn("unknown channel, dropping", "channel", data.GetChannel())

			continue
		}
	}
}
