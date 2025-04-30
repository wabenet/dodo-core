package runtime

import (
	"errors"
	"fmt"
	"strings"

	api "github.com/wabenet/dodo-core/api/runtime/v1alpha2"
)

var ErrPortFormat = errors.New("invalid publish format")

type PortConfig []PortBinding

type PortBinding struct {
	HostPort      string
	ContainerPort string
	Protocol      string
	HostIP        string
}

func MergePortConfig(first, second PortConfig) PortConfig {
	return append(first, second...)
}

func PortConfigFromProto(p []*api.PortBinding) PortConfig {
	out := PortConfig{}

	for _, port := range p {
		out = append(out, PortBindingFromProto(port))
	}

	return out
}

func (p PortConfig) ToProto() []*api.PortBinding {
	out := []*api.PortBinding{}

	for _, port := range p {
		out = append(out, port.ToProto())
	}

	return out
}

func PortBindingFromProto(p *api.PortBinding) PortBinding {
	return PortBinding{
		HostPort:      p.GetHostPort(),
		ContainerPort: p.GetContainerPort(),
		Protocol:      p.GetProtocol(),
		HostIP:        p.GetHostIp(),
	}
}

func (p PortBinding) ToProto() *api.PortBinding {
	return &api.PortBinding{
		HostPort:      p.HostPort,
		ContainerPort: p.ContainerPort,
		Protocol:      p.Protocol,
		HostIp:        p.HostIP,
	}
}

func PortBindingFromSpec(spec string) (PortBinding, error) {
	port := PortBinding{}

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return port, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	case 1:
		port.ContainerPort = values[0]
	case 2:
		port.HostPort = values[0]
		port.ContainerPort = values[1]
	case 3:
		port.HostIP = values[0]
		port.HostPort = values[1]
		port.ContainerPort = values[2]
	default:
		return port, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	switch values := strings.SplitN(port.HostPort, "/", 2); len(values) {
	case 1:
		port.HostPort = values[0]
	case 2:
		port.HostPort = values[0]
		port.Protocol = values[1]
	default:
		return port, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	return port, nil
}
