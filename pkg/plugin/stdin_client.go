package plugin

import (
	"context"
	"fmt"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
	"github.com/wabenet/dodo-core/pkg/ioutil"
	"google.golang.org/grpc"
)

type InputStreamingClient struct {
	streamInputClient api.InputStreamingPluginClient
	stdin             *grpcutil.StreamInputClient
}

func NewInputStreamingClient(conn grpc.ClientConnInterface) *InputStreamingClient {
	return &InputStreamingClient{
		streamInputClient: api.NewInputStreamingPluginClient(conn),
		stdin:             grpcutil.NewStreamInputClient(),
	}
}

func (c *InputStreamingClient) copyInputClientToStdin(containerID string, stdin io.Reader) error {
	inputClient, err := c.streamInputClient.StreamInput(context.Background())
	if err != nil {
		return fmt.Errorf("could not stream runtime input: %w", err)
	}

	req := &api.StreamInputRequest{}
	initial := &api.InitialStreamInputRequest{}

	initial.SetId(containerID)
	req.SetInitialRequest(initial)

	if err := inputClient.Send(req); err != nil {
		return fmt.Errorf("could not stream runtime input: %w", err)
	}

	if err := c.stdin.StreamInput(&streamInputClient{client: inputClient}, stdin); err != nil {
		return fmt.Errorf("could not stream runtime input: %w", err)
	}

	return nil
}

type streamInputClient struct {
	client api.InputStreamingPlugin_StreamInputClient
}

func (s *streamInputClient) Send(data *api.SubsequentStreamInputRequest) error {
	req := &api.StreamInputRequest{}

	req.SetInputData(data)

	if err := s.client.Send(req); err != nil {
		return fmt.Errorf("error wrapping Send call: %w", err)
	}

	return nil
}

func (s *streamInputClient) CloseAndRecv() (*empty.Empty, error) {
	e, err := s.client.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("error wrapping CloseAndRecv call: %w", err)
	}

	return e, nil
}

type ClientInputStream struct {
	Copy  func() error
	Close func()
}

func (c *InputStreamingClient) PrepareStream(streamID string, stdin io.Reader) (ClientInputStream, error) {
	inContext, inCancel := context.WithCancel(context.Background())
	inReader := ioutil.NewCancelableReader(inContext, stdin)

	return ClientInputStream{
		Copy: func() error {
			return c.copyInputClientToStdin(streamID, inReader)
		},
		Close: inCancel,
	}, nil
}
