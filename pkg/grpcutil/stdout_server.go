package grpcutil

import (
	"bufio"
	"context"
	"fmt"
	"io"

	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"golang.org/x/sync/errgroup"
)

type StreamOutputServer[R streamOutputResponse] struct {
	stdoutCh   chan []byte
	stderrCh   chan []byte
	outputDone chan error
}

type grpcOutputServer[R streamOutputResponse] interface {
	Send(resp R) error
	Context() context.Context
}

func NewStreamOutputServer[R streamOutputResponse]() *StreamOutputServer[R] {
	return &StreamOutputServer[R]{
		stdoutCh:   make(chan []byte),
		stderrCh:   make(chan []byte),
		outputDone: make(chan error, 1),
	}
}

func (s *StreamOutputServer[R]) ReadFrom(stdout, stderr io.Reader) error {
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

func (s *StreamOutputServer[R]) SendTo(srv grpcOutputServer[R]) error {
	var data api.OutputData

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

			data.SetData(d)
			data.SetChannel(api.OutputData_CHANNEL_STDOUT)

		case d, ok := <-s.stderrCh:
			if !ok {
				s.stderrCh = nil

				continue
			}

			data.SetData(d)
			data.SetChannel(api.OutputData_CHANNEL_STDERR)

		case <-srv.Context().Done():
			return nil
		}

		if len(data.GetData()) == 0 {
			continue
		}

		var resp R

		resp.SetOutputData(&data)

		if err := srv.Send(resp); err != nil {
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
