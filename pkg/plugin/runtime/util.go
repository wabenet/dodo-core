package runtime

import (
	log "github.com/hashicorp/go-hclog"
	api "github.com/wabenet/dodo-core/api/v1alpha4"
	"github.com/wabenet/dodo-core/pkg/plugin"
)

func GetByName(m plugin.Manager, name string) (ContainerRuntime, error) {
	for n, p := range m.GetPlugins(Type.String()) {
		if name != "" && name != n {
			continue
		}

		if rt, ok := p.(ContainerRuntime); ok {
			log.L().Info("using runtime", "name", n)

			return rt, nil
		}
	}

	return nil, plugin.NotFoundError{
		Plugin: &api.PluginName{Type: Type.String(), Name: name},
	}
}
