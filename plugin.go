package plugin

import (
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
	"github.com/wabenet/dodo-core/pkg/plugin/command"
	"github.com/wabenet/dodo-core/pkg/plugin/configuration"
)

func IncludeMe(m plugin.Manager) {
	m.RegisterPluginTypes(command.Type, configuration.Type, builder.Type)
}
