package builder_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	api "github.com/wabenet/dodo-core/api/v1alpha3"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
	plugintest "github.com/wabenet/dodo-core/pkg/plugin/test"
)

var _ builder.ImageBuilder = &DummyBuilder{}

type DummyBuilder struct{}

func (b *DummyBuilder) Type() dodo.Type {
	return builder.Type
}

func (b *DummyBuilder) PluginInfo() *api.PluginInfo {
	return &api.PluginInfo{
		Name: &api.PluginName{Type: builder.Type.String(), Name: "dummy"},
	}
}

func (b *DummyBuilder) Init() (dodo.PluginConfig, error) {
	return map[string]string{"testkey": "testvalue"}, nil
}

func (b *DummyBuilder) CreateImage(config *api.BuildInfo, stream *dodo.StreamConfig) (string, error) {
	stream.Stdout.Write([]byte("This goes to stdout"))
	stream.Stderr.Write([]byte("This goes to stderr"))

	return config.ImageName, nil
}

func TestCreateImage(t *testing.T) {
	c, cleanup, err := plugintest.GRPCWrapPlugin(builder.Type, &DummyBuilder{})
	assert.Nil(t, err)

	defer cleanup()

	b, ok := c.(builder.ImageBuilder)
	assert.True(t, ok)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	imageID, err := b.CreateImage(
		&api.BuildInfo{ImageName: "testimage"},
		&dodo.StreamConfig{
			Stdout: stdout, Stderr: stderr,
			TerminalHeight: 1, TerminalWidth: 1,
		},
	)
	assert.Nil(t, err)

	assert.Equal(t, "testimage", imageID)
	assert.Equal(t, "This goes to stdout", stdout.String())
	assert.Equal(t, "This goes to stderr", stderr.String())
}
