package builder_test

import (
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
)

var _ builder.ImageBuilder = &DummyBuilder{}

type DummyBuilder struct{}

func (b *DummyBuilder) Type() plugin.Type {
	return builder.Type
}

func (b *DummyBuilder) Metadata() plugin.Metadata {
	return plugin.NewMetadata(builder.Type, "dummy")
}

func (b *DummyBuilder) Init() (plugin.Config, error) {
	return map[string]string{"testkey": "testvalue"}, nil
}

func (*DummyBuilder) Cleanup() {}

func (b *DummyBuilder) CreateImage(config builder.BuildConfig, stream *plugin.StreamConfig) (string, error) {
	stream.Stdout.Write([]byte("This goes to stdout"))
	stream.Stderr.Write([]byte("This goes to stderr"))

	return config.ImageName, nil
}
