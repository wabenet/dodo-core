package plugin

import (
	"fmt"
	"io"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
)

type InputStreamingServer struct {
	stdin sync.Map
}

func (s *InputStreamingServer) Reset() {
	s.stdin = sync.Map{}
}

func (s *InputStreamingServer) stdinServer(streamID string) (*grpcutil.StreamInputServer, error) {
	inputServer, _ := s.stdin.LoadOrStore(streamID, grpcutil.NewStreamInputServer())

	result, ok := inputServer.(*grpcutil.StreamInputServer)
	if !ok {
		return nil, ErrUnexpectedMapType
	}

	return result, nil
}

func (s *InputStreamingServer) StreamInput(srv api.InputStreamingPlugin_StreamInputServer) error {
	req, err := srv.Recv()
	if err != nil {
		return fmt.Errorf("error during input stream: %w", err)
	}

	id := req.GetInitialRequest().GetId()

	inputServer, err := s.stdinServer(id)
	if err != nil {
		return fmt.Errorf("could not find stream input server: %w", err)
	}

	if err := inputServer.ReceiveFrom(&streamInputServer{server: srv}); err != nil {
		return fmt.Errorf("error during input stream: %w", err)
	}

	return nil
}

func copyInputServerToStdin(inputServer *grpcutil.StreamInputServer, stdin io.WriteCloser) error {
	if err := inputServer.WriteTo(stdin); err != nil {
		return fmt.Errorf("error writing input stream: %w", err)
	}

	return nil
}

type ServerInputStream struct {
	Stdin io.Reader
	Copy  func() error
	Close func()
}

func (s *InputStreamingServer) PrepareStream(streamID string) (ServerInputStream, error) {
	inReader, inWriter := io.Pipe()

	inputServer, err := s.stdinServer(streamID)
	if err != nil {
		return ServerInputStream{}, fmt.Errorf("could not find stream input server: %w", err)
	}

	return ServerInputStream{
		Stdin: inReader,
		Copy: func() error {
			return copyInputServerToStdin(inputServer, inWriter)
		},
		Close: func() {
			inputServer.Close()
		},
	}, nil
}

type streamInputServer struct {
	server api.InputStreamingPlugin_StreamInputServer
}

func (s *streamInputServer) Recv() (*api.SubsequentStreamInputRequest, error) {
	d, err := s.server.Recv()
	if err != nil {
		return nil, fmt.Errorf("error wrapping Recv call: %w", err)
	}

	return d.GetInputData(), nil
}

func (s *streamInputServer) SendAndClose(e *empty.Empty) error {
	if err := s.server.SendAndClose(e); err != nil {
		return fmt.Errorf("error wrapping SendAndClose call: %w", err)
	}

	return nil
}
