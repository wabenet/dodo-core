package builder_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	core "github.com/wabenet/dodo-core/api/core/v1alpha6"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
	plugintest "github.com/wabenet/dodo-core/pkg/plugin/test"
)

func TestCreateImage(t *testing.T) {
	t.Parallel()

	c, cleanup, err := plugintest.GRPCWrapPlugin(builder.Type, &DummyBuilder{})
	assert.Nil(t, err)

	defer cleanup()

	b, ok := c.(builder.ImageBuilder)
	assert.True(t, ok)

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	imageID, err := b.CreateImage(
		&core.BuildInfo{ImageName: "testimage"},
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
