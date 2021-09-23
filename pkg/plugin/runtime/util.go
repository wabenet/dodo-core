package runtime

import (
	"fmt"

	"github.com/dodo-cli/dodo-core/pkg/plugin"
	log "github.com/hashicorp/go-hclog"
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

	return nil, fmt.Errorf("could not find container runtime: %w", plugin.ErrNoValidPluginFound)
}
