package plugin_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

const (
	inMessage  = "Hello World!"
	outMessage = "Hello back :)"
	errMessage = "ohmygosh no"
)

func TestStdio(t *testing.T) {
	t.Parallel()

	server, err := plugin.NewStdioServer()
	assert.Nil(t, err)

	client, err := plugin.NewStdioClient(server.Endpoint())
	assert.Nil(t, err)

	stdin := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		// Avoiding race condition by waiting, what a great idea :/
		// For context: The copy for stdin is canceled once the copy of
                // stdout is done. This is completely expected behaviour, but
                // fails here in the test because we want to validate that the
                // stdin is also correct. Probably the test should be written
                // in a different way, but for now we are happy if we validate
                // that the production code works as intended.
		time.Sleep(100 * time.Millisecond)

		return server.Copy(
			stdin,
			bytes.NewBuffer([]byte(outMessage)),
			bytes.NewBuffer([]byte(errMessage)),
		)
	})

	eg.Go(func() error {
		return client.Copy(
			bytes.NewBuffer([]byte(inMessage)),
			stdout,
			stderr,
		)
	})

	err = eg.Wait()

	assert.Nil(t, err)
	assert.Equal(t, inMessage, stdin.String())
	assert.Equal(t, outMessage, stdout.String())
	assert.Equal(t, errMessage, stderr.String())
}

func TestOutputOnly(t *testing.T) {
	t.Parallel()

	server, err := plugin.NewStdioServer()
	assert.Nil(t, err)

	client, err := plugin.NewStdioClient(server.Endpoint())
	assert.Nil(t, err)

	stdin := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return server.Copy(
			stdin,
			bytes.NewBuffer([]byte(outMessage)),
			bytes.NewBuffer([]byte(errMessage)),
		)
	})

	eg.Go(func() error {
		return client.Copy(
			bytes.NewReader(nil),
			stdout,
			stderr,
		)
	})

	err = eg.Wait()

	assert.Nil(t, err)
	assert.Equal(t, outMessage, stdout.String())
	assert.Equal(t, errMessage, stderr.String())
}

func TestMuxDemux(t *testing.T) {
	t.Parallel()

	reader, writer := io.Pipe()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	muxer := plugin.NewMuxCopier(
		bytes.NewBuffer([]byte(outMessage)),
		bytes.NewBuffer([]byte(errMessage)),
		writer,
	)

	demuxer := plugin.NewDemuxCopier(
		reader,
		stdout,
		stderr,
	)

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		defer writer.Close()

		return muxer.Copy()
	})

	eg.Go(func() error {
		return demuxer.Copy()
	})

	err := eg.Wait()

	assert.Nil(t, err)
	assert.Equal(t, outMessage, stdout.String())
	assert.Equal(t, errMessage, stderr.String())
}
