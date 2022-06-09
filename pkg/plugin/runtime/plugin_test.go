package runtime_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	api "github.com/wabenet/dodo-core/api/v1alpha3"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/runtime"
	plugintest "github.com/wabenet/dodo-core/pkg/plugin/test"
)

var _ runtime.ContainerRuntime = &DummyRuntime{}

type DummyRuntime struct{}

func (r *DummyRuntime) Type() dodo.Type {
	return runtime.Type
}

func (r *DummyRuntime) PluginInfo() *api.PluginInfo {
	return &api.PluginInfo{
		Name: &api.PluginName{Type: runtime.Type.String(), Name: "dummy"},
	}
}

func (r *DummyRuntime) Init() (dodo.PluginConfig, error) {
	return map[string]string{"testkey": "testvalue"}, nil
}

func (r *DummyRuntime) ResolveImage(_ string) (string, error) {
	return "", nil
}

func (r *DummyRuntime) CreateContainer(_ *api.Backdrop, _, _ bool) (string, error) {
	return "", nil
}

func (r *DummyRuntime) StartContainer(_ string) error {
	return nil
}

func (r *DummyRuntime) DeleteContainer(_ string) error {
	return nil
}

func (r *DummyRuntime) ResizeContainer(_ string, _, _ uint32) error {
	return nil
}

func (r *DummyRuntime) KillContainer(_ string, _ os.Signal) error {
	return nil
}

func (r *DummyRuntime) StreamContainer(_ string, stream *dodo.StreamConfig) (*runtime.Result, error) {
	stream.Stderr.Write([]byte("This goes to stderr"))

	io.Copy(stream.Stdout, stream.Stdin)

	return &runtime.Result{ExitCode: 1}, nil
}

func TestStreamContainer(t *testing.T) {
	t.Parallel()

	c, cleanup, err := plugintest.GRPCWrapPlugin(runtime.Type, &DummyRuntime{})
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

	assert.Equal(t, 1, result.ExitCode)
	assert.Equal(t, "This is stdin", stdout.String())
	assert.Equal(t, "This goes to stderr", stderr.String())
}
