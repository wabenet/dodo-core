package plugin

import (
	"github.com/dodo-cli/dodo-core/pkg/command/build"
	"github.com/dodo-cli/dodo-core/pkg/command/run"
	dodo "github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/builder"
	"github.com/dodo-cli/dodo-core/pkg/plugin/command"
	"github.com/dodo-cli/dodo-core/pkg/plugin/configuration"
	"github.com/dodo-cli/dodo-core/pkg/plugin/runtime"
)

func IncludeMe() {
	dodo.RegisterPluginTypes(command.Type, configuration.Type, runtime.Type, builder.Type)
	dodo.IncludePlugins(run.New(), build.New())
}
