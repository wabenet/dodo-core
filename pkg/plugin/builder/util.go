package builder

import (
	log "github.com/hashicorp/go-hclog"
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

	return nil, plugin.NewNotFoundError(Type, name)
}
