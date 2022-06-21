package runtime_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/runtime"
	plugintest "github.com/wabenet/dodo-core/pkg/plugin/test"
)

func TestStreamStdout(t *testing.T) {
	t.Parallel()

	c, cleanup, err := plugintest.GRPCWrapPlugin(runtime.Type, &DummyRuntime{})
	assert.Nil(t, err)

	defer cleanup()

	r, ok := c.(runtime.ContainerRuntime)
	assert.True(t, ok)

	stdin, _ := io.Pipe()
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	result, err := r.StreamContainer(
		"",
		&dodo.StreamConfig{
			Stdin: stdin, Stdout: stdout, Stderr: stderr,
			TerminalHeight: 1, TerminalWidth: 1,
		},
	)
	assert.Nil(t, err)

	assert.Equal(t, 0, result.ExitCode)
	assert.Equal(t, "This goes to stdout", stdout.String())
}

func TestError(t *testing.T) {
	t.Parallel()

	c, cleanup, err := plugintest.GRPCWrapPlugin(runtime.Type, &ErrorRuntime{})
	assert.Nil(t, err)

	defer cleanup()

	r, ok := c.(runtime.ContainerRuntime)
	assert.True(t, ok)

	stdin, _ := io.Pipe()
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	result, err := r.StreamContainer(
		"",
		&dodo.StreamConfig{
			Stdin: stdin, Stdout: stdout, Stderr: stderr,
			TerminalHeight: 1, TerminalWidth: 1,
		},
	)
	assert.Nil(t, err)

	assert.Equal(t, 1, result.ExitCode)
	assert.Equal(t, "This goes to stderr", stderr.String())
}

func TestEcho(t *testing.T) {
	t.Parallel()

	c, cleanup, err := plugintest.GRPCWrapPlugin(runtime.Type, &EchoRuntime{})
	assert.Nil(t, err)

	defer cleanup()

	r, ok := c.(runtime.ContainerRuntime)
	assert.True(t, ok)

	stdin := strings.NewReader("This is stdin")
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	result, err := r.StreamContainer(
		"",
		&dodo.StreamConfig{
			Stdin: stdin, Stdout: stdout, Stderr: stderr,
			TerminalHeight: 1, TerminalWidth: 1,
		},
	)
	assert.Nil(t, err)

	assert.Equal(t, 0, result.ExitCode)
	assert.Equal(t, "This is stdin", stdout.String())
}
