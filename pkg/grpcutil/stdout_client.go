package grpcutil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	log "github.com/hashicorp/go-hclog"
	core "github.com/wabenet/dodo-core/api/core/v1alpha6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StreamOutputClient struct{}

type grpcOutputClient interface {
	Recv() (*core.OutputData, error)
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
		case core.OutputData_STDOUT:
			if _, err := io.Copy(stdout, bytes.NewReader(data.GetData())); err != nil {
				log.L().Error("failed to copy all bytes", "err", err)
			}

		case core.OutputData_STDERR:
			if _, err := io.Copy(stderr, bytes.NewReader(data.GetData())); err != nil {
				log.L().Error("failed to copy all bytes", "err", err)
			}

		case core.OutputData_INVALID:
			log.L().Warn("unknown channel, dropping", "channel", data.GetChannel())

			continue
		}
	}
}
