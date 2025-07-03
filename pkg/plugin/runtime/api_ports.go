package runtime

import (
	"errors"
	"fmt"
	"strings"

	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/runtime/v1alpha2"
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
	out := &api.PortBinding{}

	out.SetHostPort(p.HostPort)
	out.SetContainerPort(p.ContainerPort)
	out.SetProtocol(p.Protocol)
	out.SetHostIp(p.HostIP)

	return out
}

func PortBindingFromSpec(spec string) (PortBinding, error) {
	out := PortBinding{}

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return out, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	case 1:
		out.ContainerPort = values[0]
	case 2:
		out.HostPort = values[0]
		out.ContainerPort = values[1]
	case 3:
		out.HostIP = values[0]
		out.HostPort = values[1]
		out.ContainerPort = values[2]
	default:
		return out, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	switch values := strings.SplitN(out.HostPort, "/", 2); len(values) {
	case 1:
		out.HostPort = values[0]
	case 2:
		out.HostPort = values[0]
		out.Protocol = values[1]
	default:
		return out, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	return out, nil
}
