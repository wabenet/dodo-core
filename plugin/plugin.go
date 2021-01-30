package plugin

import (
	dodo "github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/command"
	"github.com/dodo-cli/dodo-core/pkg/plugin/configuration"
	"github.com/dodo-cli/dodo-core/pkg/run"
)

func IncludeMe() {
	dodo.RegisterPluginTypes(command.Type, configuration.Type)
	dodo.IncludePlugins(&run.Command{})
}
