package plugin_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
)

func TestResolveDependencies(t *testing.T) {
	t.Parallel()

	plugins := dodo.ResolveDependencies(populatePluginMap())
	assert.Equal(t, 2, len(plugins))

	resultA := plugins[0]
	assert.Equal(t, pluginAImpl, resultA)

	resultB := plugins[1]
	assert.Equal(t, pluginBImpl, resultB)
}
