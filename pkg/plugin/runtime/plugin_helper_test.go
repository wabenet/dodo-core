package runtime_test

import (
	"io"
	"os"

	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/runtime"
)

var _ runtime.ContainerRuntime = &DummyRuntime{}

type DummyRuntime struct{}

func (r *DummyRuntime) Type() plugin.Type {
	return runtime.Type
}

func (r *DummyRuntime) Metadata() plugin.Metadata {
	return plugin.NewMetadata(runtime.Type, "dummy")
}

func (r *DummyRuntime) Init() (plugin.Config, error) {
	return map[string]string{}, nil
}

func (*DummyRuntime) Cleanup() {}

func (r *DummyRuntime) ResolveImage(_ string) (string, error) {
	return "", nil
}

func (r *DummyRuntime) CreateContainer(_ runtime.ContainerConfig) (string, error) {
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

func (r *DummyRuntime) StreamContainer(_ string, stream *plugin.StreamConfig) (*runtime.Result, error) {
	stream.Stdout.Write([]byte("This goes to stdout"))

	return &runtime.Result{ExitCode: 0}, nil
}

func (r *DummyRuntime) CreateVolume(_ string) error {
	return nil
}

func (r *DummyRuntime) DeleteVolume(_ string) error {
	return nil
}

func (r *DummyRuntime) WriteFile(_, _ string, _ []byte) error {
	return nil
}

var _ runtime.ContainerRuntime = &ErrorRuntime{}

type ErrorRuntime struct{}

func (r *ErrorRuntime) Type() plugin.Type {
	return runtime.Type
}

func (r *ErrorRuntime) Metadata() plugin.Metadata {
	return plugin.NewMetadata(runtime.Type, "error")
}

func (r *ErrorRuntime) Init() (plugin.Config, error) {
	return map[string]string{}, nil
}

func (*ErrorRuntime) Cleanup() {}

func (r *ErrorRuntime) ResolveImage(_ string) (string, error) {
	return "", nil
}

func (r *ErrorRuntime) CreateContainer(_ runtime.ContainerConfig) (string, error) {
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

func (r *ErrorRuntime) StreamContainer(_ string, stream *plugin.StreamConfig) (*runtime.Result, error) {
	stream.Stdout.Write([]byte("This goes to stdout"))
	stream.Stderr.Write([]byte("This goes to stderr"))

	return &runtime.Result{ExitCode: 1}, nil
}

func (r *ErrorRuntime) CreateVolume(_ string) error {
	return nil
}

func (r *ErrorRuntime) DeleteVolume(_ string) error {
	return nil
}

func (r *ErrorRuntime) WriteFile(_, _ string, _ []byte) error {
	return nil
}

var _ runtime.ContainerRuntime = &EchoRuntime{}

type EchoRuntime struct{}

func (r *EchoRuntime) Type() plugin.Type {
	return runtime.Type
}

func (r *EchoRuntime) Metadata() plugin.Metadata {
	return plugin.NewMetadata(runtime.Type, "echo")
}

func (r *EchoRuntime) Init() (plugin.Config, error) {
	return map[string]string{}, nil
}

func (*EchoRuntime) Cleanup() {}

func (r *EchoRuntime) ResolveImage(_ string) (string, error) {
	return "", nil
}

func (r *EchoRuntime) CreateContainer(_ runtime.ContainerConfig) (string, error) {
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

func (r *EchoRuntime) StreamContainer(_ string, stream *plugin.StreamConfig) (*runtime.Result, error) {
	io.Copy(stream.Stdout, stream.Stdin)

	return &runtime.Result{ExitCode: 0}, nil
}

func (r *EchoRuntime) CreateVolume(_ string) error {
	return nil
}

func (r *EchoRuntime) DeleteVolume(_ string) error {
	return nil
}

func (r *EchoRuntime) WriteFile(_, _ string, _ []byte) error {
	return nil
}
