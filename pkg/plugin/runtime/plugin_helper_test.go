package runtime_test

import (
	"io"
	"os"

	api "github.com/wabenet/dodo-core/api/v1alpha3"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/runtime"
)

var _ runtime.ContainerRuntime = &DummyRuntime{}

type DummyRuntime struct{}

func (r *DummyRuntime) Type() dodo.Type {
	return runtime.Type
}

func (r *DummyRuntime) PluginInfo() *api.PluginInfo {
	return &api.PluginInfo{
		Name: &api.PluginName{Type: runtime.Type.String(), Name: "dummy"},
	}
}

func (r *DummyRuntime) Init() (dodo.PluginConfig, error) {
	return map[string]string{}, nil
}

func (r *DummyRuntime) ResolveImage(_ string) (string, error) {
	return "", nil
}

func (r *DummyRuntime) CreateContainer(_ *api.Backdrop, _, _ bool) (string, error) {
	return "", nil
}

func (r *DummyRuntime) StartContainer(_ string) error {
	return nil
}

func (r *DummyRuntime) DeleteContainer(_ string) error {
	return nil
}

func (r *DummyRuntime) ResizeContainer(_ string, _, _ uint32) error {
	return nil
}

func (r *DummyRuntime) KillContainer(_ string, _ os.Signal) error {
	return nil
}

func (r *DummyRuntime) StreamContainer(_ string, stream *dodo.StreamConfig) (*runtime.Result, error) {
	stream.Stdout.Write([]byte("This goes to stdout"))

	return &runtime.Result{ExitCode: 0}, nil
}

var _ runtime.ContainerRuntime = &ErrorRuntime{}

type ErrorRuntime struct{}

func (r *ErrorRuntime) Type() dodo.Type {
	return runtime.Type
}

func (r *ErrorRuntime) PluginInfo() *api.PluginInfo {
	return &api.PluginInfo{
		Name: &api.PluginName{Type: runtime.Type.String(), Name: "error"},
	}
}

func (r *ErrorRuntime) Init() (dodo.PluginConfig, error) {
	return map[string]string{}, nil
}

func (r *ErrorRuntime) ResolveImage(_ string) (string, error) {
	return "", nil
}

func (r *ErrorRuntime) CreateContainer(_ *api.Backdrop, _, _ bool) (string, error) {
	return "", nil
}

func (r *ErrorRuntime) StartContainer(_ string) error {
	return nil
}

func (r *ErrorRuntime) DeleteContainer(_ string) error {
	return nil
}

func (r *ErrorRuntime) ResizeContainer(_ string, _, _ uint32) error {
	return nil
}

func (r *ErrorRuntime) KillContainer(_ string, _ os.Signal) error {
	return nil
}

func (r *ErrorRuntime) StreamContainer(_ string, stream *dodo.StreamConfig) (*runtime.Result, error) {
	stream.Stdout.Write([]byte("This goes to stdout"))
	stream.Stderr.Write([]byte("This goes to stderr"))

	return &runtime.Result{ExitCode: 1}, nil
}

var _ runtime.ContainerRuntime = &EchoRuntime{}

type EchoRuntime struct{}

func (r *EchoRuntime) Type() dodo.Type {
	return runtime.Type
}

func (r *EchoRuntime) PluginInfo() *api.PluginInfo {
	return &api.PluginInfo{
		Name: &api.PluginName{Type: runtime.Type.String(), Name: "echo"},
	}
}

func (r *EchoRuntime) Init() (dodo.PluginConfig, error) {
	return map[string]string{}, nil
}

func (r *EchoRuntime) ResolveImage(_ string) (string, error) {
	return "", nil
}

func (r *EchoRuntime) CreateContainer(_ *api.Backdrop, _, _ bool) (string, error) {
	return "", nil
}

func (r *EchoRuntime) StartContainer(_ string) error {
	return nil
}

func (r *EchoRuntime) DeleteContainer(_ string) error {
	return nil
}

func (r *EchoRuntime) ResizeContainer(_ string, _, _ uint32) error {
	return nil
}

func (r *EchoRuntime) KillContainer(_ string, _ os.Signal) error {
	return nil
}

func (r *EchoRuntime) StreamContainer(_ string, stream *dodo.StreamConfig) (*runtime.Result, error) {
	io.Copy(stream.Stdout, stream.Stdin)

	return &runtime.Result{ExitCode: 0}, nil
}
