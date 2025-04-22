package config

import (
	"fmt"
	"os"
	"strings"

	build "github.com/wabenet/dodo-core/api/build/v1alpha2"
	runtime "github.com/wabenet/dodo-core/api/runtime/v1alpha2"
)

type FormatError string

func (e FormatError) Error() string {
	return string(e)
}

const (
	ErrArgumentFormat    FormatError = "invalid argument format"
	ErrDeviceFormat      FormatError = "invalid device format"
	ErrVolumeFormat      FormatError = "invalid volume format"
	ErrEnvironmentFormat FormatError = "invalid environment format"
	ErrPortFormat        FormatError = "invalid publish format"
	ErrSSHAgentFormat    FormatError = "invalid ssh agent format"
)

func ParseBuildArgument(spec string) (*build.BuildArgument, error) {
	arg := &build.BuildArgument{}

	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 0:
		return nil, fmt.Errorf("%s: %w", spec, ErrArgumentFormat)
	case 1:
		arg.Key = values[0]
		arg.Value = os.Getenv(values[0])
	case 2:
		arg.Key = values[0]
		arg.Value = values[1]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrArgumentFormat)
	}

	return arg, nil
}

func ParseDeviceMount(spec string) (*runtime.Mount, error) {
	dev := &runtime.DeviceMount{}

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return nil, fmt.Errorf("%s: %w", spec, ErrDeviceFormat)
	case 1:
		dev.HostPath = values[0]
	case 2:
		dev.HostPath = values[0]
		dev.ContainerPath = values[1]
	case 3:
		dev.HostPath = values[0]
		dev.ContainerPath = values[1]
		dev.Permissions = values[2]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrDeviceFormat)
	}

	return &runtime.Mount{
		Type: &runtime.Mount_Device{Device: dev},
	}, nil
}

func ParseBindMount(spec string) (*runtime.Mount, error) {
	var source, target string
	var readonly bool

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return nil, fmt.Errorf("%s: %w", spec, ErrVolumeFormat)
	case 1:
		source = values[0]
	case 2:
		source = values[0]
		target = values[1]
	case 3:
		source = values[0]
		target = values[1]
		readonly = (values[2] == "ro")
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrVolumeFormat)
	}

	// If source is an absolute path, assume bind mount, otherwise assume a volume
	if source[0] == '/' {
		return &runtime.Mount{
			Type: &runtime.Mount_Bind{Bind: &runtime.BindMount{
				HostPath:      source,
				ContainerPath: target,
				Readonly:      readonly,
			}},
		}, nil
	} else {
		return &runtime.Mount{
			Type: &runtime.Mount_Volume{Volume: &runtime.VolumeMount{
				VolumeName:    source,
				ContainerPath: target,
				Readonly:      readonly,
			}},
		}, nil
	}
}

func ParseEnvironmentVariable(spec string) (*runtime.EnvironmentVariable, error) {
	env := &runtime.EnvironmentVariable{}

	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 0:
		return nil, fmt.Errorf("%s: %w", spec, ErrEnvironmentFormat)
	case 1:
		env.Key = values[0]
		env.Value = os.Getenv(values[0])
	case 2:
		env.Key = values[0]
		env.Value = values[1]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrEnvironmentFormat)
	}

	return env, nil
}

func ParsePortBinding(spec string) (*runtime.PortBinding, error) {
	port := &runtime.PortBinding{}

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return nil, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	case 1:
		port.ContainerPort = values[0]
	case 2:
		port.HostPort = values[0]
		port.ContainerPort = values[1]
	case 3:
		port.HostIp = values[0]
		port.HostPort = values[1]
		port.ContainerPort = values[2]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	switch values := strings.SplitN(port.GetHostPort(), "/", 2); len(values) {
	case 1:
		port.HostPort = values[0]
	case 2:
		port.HostPort = values[0]
		port.Protocol = values[1]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	return port, nil
}

func ParseSSHAgent(spec string) (*build.SshAgent, error) {
	agent := &build.SshAgent{}

	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 2:
		agent.Id = values[0]
		agent.IdentityFile = values[1]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrSSHAgentFormat)
	}

	return agent, nil
}
