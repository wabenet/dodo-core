package stdio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InputStreamingServer struct {
	stdin sync.Map
}

func (s *InputStreamingServer) Reset() {
	s.stdin = sync.Map{}
}

func (s *InputStreamingServer) stdinServer(streamID string) (*StreamInputServer, error) {
	inputServer, _ := s.stdin.LoadOrStore(streamID, NewStreamInputServer())

	result, ok := inputServer.(*StreamInputServer)
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

func copyInputServerToStdin(inputServer *StreamInputServer, stdin io.WriteCloser) error {
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

type StreamInputServer struct {
	stdinCh     chan []byte
	inputDone   chan error
	stdinCloser sync.Once
}

type grpcInputServer interface {
	Recv() (*api.SubsequentStreamInputRequest, error)
	SendAndClose(_ *empty.Empty) error
}

func NewStreamInputServer() *StreamInputServer {
	return &StreamInputServer{
		stdinCh:     make(chan []byte),
		inputDone:   make(chan error, 1),
		stdinCloser: sync.Once{},
	}
}

func (s *StreamInputServer) WriteTo(stdin io.WriteCloser) error {
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		defer stdin.Close()

		return copyInput(stdin, s.stdinCh)
	})

	eg.Go(func() error {
		return <-s.inputDone
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("error writing input stream: %w", err)
	}

	return nil
}

func (s *StreamInputServer) ReceiveFrom(srv grpcInputServer) error {
	defer s.Close()

	for {
		data, err := srv.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if err := srv.SendAndClose(&empty.Empty{}); err != nil {
					return fmt.Errorf("could not close input stream: %w", err)
				}

				return nil
			}

			if errors.Is(err, context.Canceled) ||
				status.Code(err) == codes.Unavailable ||
				status.Code(err) == codes.Canceled ||
				status.Code(err) == codes.Unimplemented {
				return nil
			}

			log.L().Error("error receiving data", "err", err)

			return fmt.Errorf("error receiving build input from client: %w", err)
		}

		s.stdinCh <- data.GetData()
	}
}

func (s *StreamInputServer) Close() {
	s.stdinCloser.Do(func() {
		close(s.stdinCh)
		s.inputDone <- nil
	})
}

func copyInput(dst io.Writer, src chan []byte) error {
	for data := range src {
		if len(data) == 0 {
			continue
		}

		if _, err := dst.Write(data); err != nil {
			log.L().Warn("error in stdio stream", "err", err)

			break
		}
	}

	return nil
}
