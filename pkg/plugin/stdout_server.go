package plugin

import (
	"errors"
	"fmt"
	"io"
	"sync"

	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/plugin/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/grpcutil"
)

var ErrUnexpectedMapType = errors.New("unexpected map type for stdio streaming server")

type OutputStreamingServer struct {
	stdout sync.Map
}

func (s *OutputStreamingServer) Reset() {
	s.stdout = sync.Map{}
}

func (s *OutputStreamingServer) stdoutServer(streamId string) (*grpcutil.StreamOutputServer, error) {
	outputServer, _ := s.stdout.LoadOrStore(streamId, grpcutil.NewStreamOutputServer())

	result, ok := outputServer.(*grpcutil.StreamOutputServer)
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

func copyOutputServerToStdout(outputServer *grpcutil.StreamOutputServer, stdout, stderr io.Reader) error {
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
