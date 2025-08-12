package stdio

import (
	"context"
	"errors"
	"fmt"
	"io"

	log "github.com/hashicorp/go-hclog"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/ioutil"
)

type ClientInputStream struct {
	client     api.InputStreamingPluginClient
	streamID   string
	stdin      io.Reader
	cancelFunc func()
}

func NewClientInputStream(client api.InputStreamingPluginClient, streamID string, stdin io.Reader) (ClientInputStream, error) {
	inContext, inCancel := context.WithCancel(context.Background())
	inReader := ioutil.NewCancelableReader(inContext, stdin)

	return ClientInputStream{
		client:     client,
		streamID:   streamID,
		stdin:      inReader,
		cancelFunc: inCancel,
	}, nil
}

func (s ClientInputStream) Copy() error {
	inputClient, err := s.client.StreamInput(context.Background())
	if err != nil {
		return fmt.Errorf("could not stream runtime input: %w", err)
	}

	req := &api.StreamInputRequest{}
	initial := &api.InitialStreamInputRequest{}

	initial.SetId(s.streamID)
	req.SetInitialRequest(initial)

	if err := inputClient.Send(req); err != nil {
		return fmt.Errorf("could not stream runtime input: %w", err)
	}

	data := api.SubsequentStreamInputRequest{}

	for {
		var b [1024]byte

		n, err := s.stdin.Read(b[:])

		if n > 0 {
			req := &api.StreamInputRequest{}

			data.SetData(b[:n])
			req.SetInputData(&data)

			if serr := inputClient.Send(req); serr != nil {
				return fmt.Errorf("could not send input to server: %w", serr)
			}
		}

		if errors.Is(err, io.EOF) || errors.Is(err, context.Canceled) {
			if _, serr := inputClient.CloseAndRecv(); serr != nil {
				log.L().Warn("could not close input stream", "err", serr)
			}

			return nil
		}

		if err != nil {
			return fmt.Errorf("could not read stream from client: %w", err)
		}
	}
}

func (s ClientInputStream) Close() {
	s.cancelFunc()
}
