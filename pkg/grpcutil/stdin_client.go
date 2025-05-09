package grpcutil

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	api "github.com/wabenet/dodo-core/api/plugin/v1alpha2"
)

type StreamInputClient struct{}

type grpcInputClient interface {
	Send(data *api.InputData) error
	CloseAndRecv() (*empty.Empty, error)
}

func NewStreamInputClient() *StreamInputClient {
	return &StreamInputClient{}
}

func (*StreamInputClient) StreamInput(cl grpcInputClient, stdin io.Reader) error {
	data := api.InputData{}

	for {
		var b [1024]byte

		n, err := stdin.Read(b[:])

		if n > 0 {
			data.SetData(b[:n])
			if serr := cl.Send(&data); err != nil {
				return fmt.Errorf("could not send input to server: %w", serr)
			}
		}

		if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
			if _, serr := cl.CloseAndRecv(); serr != nil {
				log.L().Warn("could not close input stream", "err", serr)
			}

			return nil
		}

		if err != nil {
			return fmt.Errorf("could not read stream from client: %w", err)
		}
	}
}
