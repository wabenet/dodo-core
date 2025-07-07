package grpcutil

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

type StreamOutputClient[R streamOutputResponse] struct{}

type streamOutputResponse interface {
	GetOutputData() *api.OutputData
	SetOutputData(data *api.OutputData)
}

type grpcOutputClient[R streamOutputResponse] interface {
	Recv() (R, error)
}

func NewStreamOutputClient[R streamOutputResponse]() *StreamOutputClient[R] {
	return &StreamOutputClient[R]{}
}

func (*StreamOutputClient[R]) StreamOutput(cl grpcOutputClient[R], stdout, stderr io.Writer) error {
	for {
		resp, err := cl.Recv()
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

		data := resp.GetOutputData()

		switch data.GetChannel() {
		case api.OutputData_CHANNEL_STDOUT:
			if _, err := io.Copy(stdout, bytes.NewReader(data.GetData())); err != nil {
				log.L().Error("failed to copy all bytes", "err", err)
			}

		case api.OutputData_CHANNEL_STDERR:
			if _, err := io.Copy(stderr, bytes.NewReader(data.GetData())); err != nil {
				log.L().Error("failed to copy all bytes", "err", err)
			}

		case api.OutputData_CHANNEL_UNSPECIFIED:
			log.L().Warn("unknown channel, dropping", "channel", data.GetChannel())

			continue
		}
	}
}
