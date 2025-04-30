package builder

import (
	"github.com/wabenet/dodo-core/pkg/plugin"
)

type ImageBuilder interface {
	plugin.Plugin

	CreateImage(config BuildConfig, streamConfig *plugin.StreamConfig) (string, error)
}
