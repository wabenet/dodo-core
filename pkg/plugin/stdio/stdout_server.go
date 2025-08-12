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

var _ api.OutputStreamingPluginServer = &OutputStreamingPluginServer{}

type OutputStreamingPluginServer struct {
	api.UnsafeOutputStreamingPluginServer

	stdout sync.Map
}

type OutputStreamingServerInstance struct {
	stdoutCh   chan []byte
	stderrCh   chan []byte
	outputDone chan error
}

type ServerOutputStream struct {
	server *OutputStreamingServerInstance

	stdoutReader io.Reader
	Stdout       io.WriteCloser
	stderrReader io.Reader
	Stderr       io.WriteCloser
}

func (s *OutputStreamingPluginServer) Reset() {
	s.stdout = sync.Map{}
}

func (s *OutputStreamingPluginServer) StreamOutput(request *api.StreamOutputRequest, srv api.OutputStreamingPlugin_StreamOutputServer) error {
	id := request.GetId()

	instance, err := s.serverInstanceForID(id)
	if err != nil {
		return fmt.Errorf("could not find stream output server: %w", err)
	}

	var response api.StreamOutputResponse

	defer func() {
		instance.outputDone <- nil
	}()

	for {
		if instance.stdoutCh == nil && instance.stderrCh == nil {
			return nil
		}

		select {
		case data, ok := <-instance.stdoutCh:
			if !ok {
				instance.stdoutCh = nil

				continue
			}

			response.SetData(data)
			response.SetChannel(api.OutputChannel_OUTPUT_CHANNEL_STDOUT)

		case data, ok := <-instance.stderrCh:
			if !ok {
				instance.stderrCh = nil

				continue
			}

			response.SetData(data)
			response.SetChannel(api.OutputChannel_OUTPUT_CHANNEL_STDERR)

		case <-srv.Context().Done():
			return nil
		}

		if len(response.GetData()) == 0 {
			continue
		}

		if err := srv.Send(&response); err != nil {
			return fmt.Errorf("error sending build output to client: %w", err)
		}
	}
}

func (s *OutputStreamingPluginServer) NewServerOutputStream(streamID string) (ServerOutputStream, error) {
	outReader, outWriter := io.Pipe()
	errReader, errWriter := io.Pipe()

	instance, err := s.serverInstanceForID(streamID)
	if err != nil {
		return ServerOutputStream{}, fmt.Errorf("could not find stream output server: %w", err)
	}

	return ServerOutputStream{
		server:       instance,
		stdoutReader: outReader,
		Stdout:       outWriter,
		stderrReader: errReader,
		Stderr:       errWriter,
	}, nil
}

func (s *OutputStreamingPluginServer) serverInstanceForID(streamID string) (*OutputStreamingServerInstance, error) {
	outputServer, _ := s.stdout.LoadOrStore(streamID, &OutputStreamingServerInstance{
		stdoutCh:   make(chan []byte),
		stderrCh:   make(chan []byte),
		outputDone: make(chan error, 1),
	})

	result, ok := outputServer.(*OutputStreamingServerInstance)
	if !ok {
		return nil, ErrUnexpectedMapType
	}

	return result, nil
}

func (s ServerOutputStream) Copy() error {
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return copyOutput(s.server.stdoutCh, s.stdoutReader)
	})

	eg.Go(func() error {
		return copyOutput(s.server.stderrCh, s.stderrReader)
	})

	eg.Go(func() error {
		return <-s.server.outputDone
	})

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("error reading output stream: %w", err)
	}

	return nil
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

func (s ServerOutputStream) Close() {
	s.Stdout.Close()
	s.Stderr.Close()
}
