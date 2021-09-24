package plugin_test

import (
	"testing"

	dodo "github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/stretchr/testify/assert"
)

func TestResolveDependencies(t *testing.T) {
	t.Parallel()

	plugins, err := dodo.ResolveDependencies(populatePluginMap())
	assert.Nil(t, err)

	assert.Equal(t, 2, len(plugins))

	resultA := plugins[0]
	assert.Equal(t, pluginAImpl, resultA)

	resultB := plugins[1]
	assert.Equal(t, pluginBImpl, resultB)
}
