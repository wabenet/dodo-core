package builder

import (
	"fmt"

	"github.com/dodo-cli/dodo-core/pkg/plugin"
	log "github.com/hashicorp/go-hclog"
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

	return nil, fmt.Errorf("could not find image builder: %w", plugin.ErrNoValidPluginFound)
}
