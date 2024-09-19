package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	api "github.com/wabenet/dodo-core/api/core/v1alpha5"
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
	var target, published, protocol, hostip string
	port := &api.PortBinding{}

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return nil, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	case 1:
		target = values[0]
	case 2:
		published = values[0]
		target = values[1]
	case 3:
		hostip = values[0]
		published = values[1]
		target = values[2]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	switch values := strings.SplitN(target, "/", 2); len(values) {
	case 1:
		target = values[0]
	case 2:
		target = values[0]
		protocol = values[1]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	if p, err := strconv.ParseInt(target, 10, 32); err != nil {
		return nil, fmt.Errorf("%s: %w", target, err)
	} else {
		port.Target = int32(p)
	}

	if p, err := strconv.ParseInt(published, 10, 32); err != nil {
		return nil, fmt.Errorf("%s: %w", published, err)
	} else {
		port.Published = int32(p)
	}

	switch protocol {
	case "tcp", "udp", "sctp":
		port.Protocol = protocol
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	port.HostIp = hostip

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
