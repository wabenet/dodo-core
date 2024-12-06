package grpcutil

import (
	"bufio"
	"context"
	"fmt"
	"io"

	core "github.com/wabenet/dodo-core/api/core/v1alpha6"
	"golang.org/x/sync/errgroup"
)

type StreamOutputServer struct {
	stdoutCh   chan []byte
	stderrCh   chan []byte
	outputDone chan error
}

type grpcOutputServer interface {
	Send(data *core.OutputData) error
	Context() context.Context
}

func NewStreamOutputServer() *StreamOutputServer {
	return &StreamOutputServer{
		stdoutCh:   make(chan []byte),
		stderrCh:   make(chan []byte),
		outputDone: make(chan error, 1),
	}
}

func (s *StreamOutputServer) ReadFrom(stdout, stderr io.Reader) error {
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return copyOutput(s.stdoutCh, stdout)
	})

	eg.Go(func() error {
		return copyOutput(s.stderrCh, stderr)
	})

	eg.Go(func() error {
		return <-s.outputDone
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("error reading output stream: %w", err)
	}

	return nil
}

func (s *StreamOutputServer) SendTo(srv grpcOutputServer) error {
	var data core.OutputData

	defer func() {
		s.outputDone <- nil
	}()

	for {
		if s.stdoutCh == nil && s.stderrCh == nil {
			return nil
		}

		select {
		case d, ok := <-s.stdoutCh:
			if !ok {
				s.stdoutCh = nil

				continue
			}

			data.Data = d
			data.Channel = core.OutputData_STDOUT

		case d, ok := <-s.stderrCh:
			if !ok {
				s.stderrCh = nil

				continue
			}

			data.Data = d
			data.Channel = core.OutputData_STDERR

		case <-srv.Context().Done():
			return nil
		}

		if len(data.GetData()) == 0 {
			continue
		}

		if err := srv.Send(&data); err != nil {
			return fmt.Errorf("error sending build output to client: %w", err)
		}
	}
}

func copyOutput(dst chan []byte, src io.Reader) error {
	defer close(dst)

	bufsrc := bufio.NewReader(src)

	for {
		var data [1024]byte

		n, err := bufsrc.Read(data[:])

		if n > 0 {
			dst <- data[:n]
		}

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return fmt.Errorf("error copying container output: %w", err)
		}
	}
}
