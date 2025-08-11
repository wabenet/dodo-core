package stdio

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"golang.org/x/sync/errgroup"
)

var ErrUnexpectedMapType = errors.New("unexpected map type for stdio streaming server")

type OutputStreamingServer struct {
	stdout sync.Map
}

func (s *OutputStreamingServer) Reset() {
	s.stdout = sync.Map{}
}

func (s *OutputStreamingServer) stdoutServer(streamId string) (*StreamOutputServer, error) {
	outputServer, _ := s.stdout.LoadOrStore(streamId, NewStreamOutputServer())

	result, ok := outputServer.(*StreamOutputServer)
	if !ok {
		return nil, ErrUnexpectedMapType
	}

	return result, nil
}

func (s *OutputStreamingServer) StreamOutput(request *api.StreamOutputRequest, srv api.OutputStreamingPlugin_StreamOutputServer) error {
	id := request.GetId()

	outputServer, err := s.stdoutServer(id)
	if err != nil {
		return fmt.Errorf("could not find stream output server: %w", err)
	}

	if err := outputServer.SendTo(srv); err != nil {
		return fmt.Errorf("error during output stream: %w", err)
	}

	return nil
}

func copyOutputServerToStdout(outputServer *StreamOutputServer, stdout, stderr io.Reader) error {
	if err := outputServer.ReadFrom(stdout, stderr); err != nil {
		return fmt.Errorf("error reading output stream: %w", err)
	}

	return nil
}

type ServerOutputStream struct {
	Stdout io.Writer
	Stderr io.Writer
	Copy   func() error
	Close  func()
}

func (s *OutputStreamingServer) PrepareStream(streamID string) (ServerOutputStream, error) {
	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()

	outputServer, err := s.stdoutServer(streamID)
	if err != nil {
		return ServerOutputStream{}, fmt.Errorf("could not find stream output server: %w", err)
	}

	return ServerOutputStream{
		Stdout: outWriter,
		Stderr: errWriter,
		Copy: func() error {
			return copyOutputServerToStdout(outputServer, outReader, errReader)
		},
		Close: func() {
			outWriter.Close()
			errWriter.Close()
		},
	}, nil
}

type StreamOutputServer struct {
	stdoutCh   chan []byte
	stderrCh   chan []byte
	outputDone chan error
}

type grpcOutputServer interface {
	Send(data *api.StreamOutputResponse) error
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
	var data api.StreamOutputResponse

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
			data.SetChannel(api.OutputChannel_OUTPUT_CHANNEL_STDOUT)

		case d, ok := <-s.stderrCh:
			if !ok {
				s.stderrCh = nil

				continue
			}

			data.SetData(d)
			data.SetChannel(api.OutputChannel_OUTPUT_CHANNEL_STDERR)

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
