package config

import (
	"fmt"
	"os"
	"strings"

	api "github.com/wabenet/dodo-core/api/core/v1alpha7"
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

func ParseBuildArgument(spec string) (*api.BuildArgument, error) {
	arg := &api.BuildArgument{}

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

func ParseDeviceMapping(spec string) (*api.DeviceMapping, error) {
	dev := &api.DeviceMapping{}

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return nil, fmt.Errorf("%s: %w", spec, ErrDeviceFormat)
	case 1:
		dev.Source = values[0]
	case 2:
		dev.Source = values[0]
		dev.Target = values[1]
	case 3:
		dev.Source = values[0]
		dev.Target = values[1]
		dev.Permissions = values[2]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrDeviceFormat)
	}

	return dev, nil
}

func ParseVolumeMount(spec string) (*api.VolumeMount, error) {
	vol := &api.VolumeMount{}

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return nil, fmt.Errorf("%s: %w", spec, ErrVolumeFormat)
	case 1:
		vol.Source = values[0]
	case 2:
		vol.Source = values[0]
		vol.Target = values[1]
	case 3:
		vol.Source = values[0]
		vol.Target = values[1]
		vol.Readonly = (values[2] == "ro")
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrVolumeFormat)
	}

	return vol, nil
}

func ParseEnvironmentVariable(spec string) (*api.EnvironmentVariable, error) {
	env := &api.EnvironmentVariable{}

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

func ParsePortBinding(spec string) (*api.PortBinding, error) {
	port := &api.PortBinding{}

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return nil, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	case 1:
		port.Target = values[0]
	case 2:
		port.Published = values[0]
		port.Target = values[1]
	case 3:
		port.HostIp = values[0]
		port.Published = values[1]
		port.Target = values[2]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	switch values := strings.SplitN(port.GetTarget(), "/", 2); len(values) {
	case 1:
		port.Target = values[0]
	case 2:
		port.Target = values[0]
		port.Protocol = values[1]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	return port, nil
}

func ParseSSHAgent(spec string) (*api.SshAgent, error) {
	agent := &api.SshAgent{}

	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 2:
		agent.Id = values[0]
		agent.IdentityFile = values[1]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrSSHAgentFormat)
	}

	return agent, nil
}
