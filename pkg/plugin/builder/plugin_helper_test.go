package builder_test

import (
	core "github.com/wabenet/dodo-core/api/core/v1alpha5"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
)

var _ builder.ImageBuilder = &DummyBuilder{}

type DummyBuilder struct{}

func (b *DummyBuilder) Type() dodo.Type {
	return builder.Type
}

func (b *DummyBuilder) PluginInfo() *core.PluginInfo {
	return &core.PluginInfo{
		Name: &core.PluginName{Type: builder.Type.String(), Name: "dummy"},
	}
}

func (b *DummyBuilder) Init() (dodo.PluginConfig, error) {
	return map[string]string{"testkey": "testvalue"}, nil
}

func (*DummyBuilder) Cleanup() {}

func (b *DummyBuilder) CreateImage(config *core.BuildInfo, stream *dodo.StreamConfig) (string, error) {
	stream.Stdout.Write([]byte("This goes to stdout"))
	stream.Stderr.Write([]byte("This goes to stderr"))

	return config.ImageName, nil
}
