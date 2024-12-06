package builder_test

import (
	core "github.com/wabenet/dodo-core/api/core/v1alpha6"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
)

var _ builder.ImageBuilder = &DummyBuilder{}

type DummyBuilder struct{}

func (b *DummyBuilder) Type() plugin.Type {
	return builder.Type
}

func (b *DummyBuilder) PluginInfo() *core.PluginInfo {
	return &core.PluginInfo{
		Name: &core.PluginName{Type: builder.Type.String(), Name: "dummy"},
	}
}

func (b *DummyBuilder) Init() (plugin.Config, error) {
	return map[string]string{"testkey": "testvalue"}, nil
}

func (*DummyBuilder) Cleanup() {}

func (b *DummyBuilder) CreateImage(config *core.BuildInfo, stream *plugin.StreamConfig) (string, error) {
	stream.Stdout.Write([]byte("This goes to stdout"))
	stream.Stderr.Write([]byte("This goes to stderr"))

	return config.GetImageName(), nil
}
