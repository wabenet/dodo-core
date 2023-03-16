package builder

import (
	log "github.com/hashicorp/go-hclog"
	core "github.com/wabenet/dodo-core/api/core/v1alpha5"
	"github.com/wabenet/dodo-core/pkg/plugin"
)

func GetByName(m plugin.Manager, name string) (ImageBuilder, error) {
	for n, p := range m.GetPlugins(Type.String()) {
		if name != "" && name != n {
			continue
		}

		if rt, ok := p.(ImageBuilder); ok {
			log.L().Info("using builder", "name", n)

			return rt, nil
		}
	}

	return nil, plugin.NotFoundError{
		Plugin: &core.PluginName{Type: Type.String(), Name: name},
	}
}
