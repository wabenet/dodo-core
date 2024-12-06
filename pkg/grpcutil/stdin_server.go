package grpcutil

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	core "github.com/wabenet/dodo-core/api/core/v1alpha6"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StreamInputServer struct {
	stdinCh     chan []byte
	inputDone   chan error
	stdinCloser sync.Once
}

type grpcInputServer interface {
	Recv() (*core.InputData, error)
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
