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

var _ api.InputStreamingPluginServer = &InputStreamingPluginServer{}

type InputStreamingPluginServer struct {
	api.UnsafeInputStreamingPluginServer

	stdin sync.Map
}

type InputStreamingServerInstance struct {
	stdinCh     chan []byte
	inputDone   chan error
	stdinCloser sync.Once
}

type ServerInputStream struct {
	server *InputStreamingServerInstance

	Stdin       io.Reader
	stdinWriter io.WriteCloser
}

func (s *InputStreamingPluginServer) Reset() {
	s.stdin = sync.Map{}
}

func (s *InputStreamingPluginServer) StreamInput(srv api.InputStreamingPlugin_StreamInputServer) error {
	req, err := srv.Recv()
	if err != nil {
		return fmt.Errorf("error during input stream: %w", err)
	}

	id := req.GetInitialRequest().GetId()

	server, err := s.serverInstanceForID(id)
	if err != nil {
		return fmt.Errorf("could not find stream input server: %w", err)
	}

	defer server.Close()

	for {
		data, err := srv.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				if err := srv.SendAndClose(&empty.Empty{}); err != nil {
					return fmt.Errorf("error wrapping SendAndClose call: %w", err)
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

		server.stdinCh <- data.GetInputData().GetData()
	}
}

func (s *InputStreamingPluginServer) NewServerInputStream(streamID string) (ServerInputStream, error) {
	inReader, inWriter := io.Pipe()

	instance, err := s.serverInstanceForID(streamID)
	if err != nil {
		return ServerInputStream{}, fmt.Errorf("could not find stream input server: %w", err)
	}

	return ServerInputStream{
		server:      instance,
		Stdin:       inReader,
		stdinWriter: inWriter,
	}, nil
}

func (s *InputStreamingPluginServer) serverInstanceForID(streamID string) (*InputStreamingServerInstance, error) {
	inputServer, _ := s.stdin.LoadOrStore(streamID, &InputStreamingServerInstance{
		stdinCh:     make(chan []byte),
		inputDone:   make(chan error, 1),
		stdinCloser: sync.Once{},
	})

	result, ok := inputServer.(*InputStreamingServerInstance)
	if !ok {
		return nil, ErrUnexpectedMapType
	}

	return result, nil
}

func (s *InputStreamingServerInstance) Close() {
	s.stdinCloser.Do(func() {
		close(s.stdinCh)

		s.inputDone <- nil
	})
}

func (s ServerInputStream) Copy() error {
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		defer s.stdinWriter.Close()

		return copyInput(s.stdinWriter, s.server.stdinCh)
	})

	eg.Go(func() error {
		return <-s.server.inputDone
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("error writing input stream: %w", err)
	}

	return nil
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

func (s ServerInputStream) Close() {
	s.server.Close()
}
