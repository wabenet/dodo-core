package runtime

import (
	"os"

	"github.com/wabenet/dodo-core/pkg/plugin"
)

type ContainerRuntime interface {
	plugin.Plugin

	ResolveImage(spec string) (string, error)
	CreateContainer(containerConfig ContainerConfig) (string, error)
	StartContainer(id string) error
	DeleteContainer(id string) error
	ResizeContainer(id string, height, width uint32) error
	KillContainer(id string, signal os.Signal) error
	StreamContainer(id string, streamConfig *plugin.StreamConfig) (*Result, error)
	CreateVolume(name string) error
	DeleteVolume(name string) error
	WriteFile(name, path string, contents []byte) error
}

type Result struct {
	ExitCode int
}
