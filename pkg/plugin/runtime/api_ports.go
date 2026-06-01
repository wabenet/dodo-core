package runtime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/runtime/v1alpha2"
)

var ErrPortFormat = errors.New("invalid publish format")

type PortConfig []PortBinding

type PortBinding struct {
	HostPort      uint16
	ContainerPort uint16
	Protocol      Protocol
	HostIP        string
}

type Protocol string

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
)

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
	out := PortBinding{
		HostPort:      uint16(p.GetHostPort()),
		ContainerPort: uint16(p.GetContainerPort()),
		HostIP:        p.GetHostIp(),
	}

	switch p.GetProtocol() {
	case api.Protocol_PROTOCOL_TCP:
		out.Protocol = TCP
	case api.Protocol_PROTOCOL_UDP:
		out.Protocol = UDP
	}

	return out
}

func (p PortBinding) ToProto() *api.PortBinding {
	out := &api.PortBinding{}

	out.SetHostPort(uint32(p.HostPort))
	out.SetContainerPort(uint32(p.ContainerPort))
	out.SetHostIp(p.HostIP)

	switch p.Protocol {
	case TCP:
		out.SetProtocol(api.Protocol_PROTOCOL_TCP)
	case UDP:
		out.SetProtocol(api.Protocol_PROTOCOL_UDP)
	default:
		out.SetProtocol(api.Protocol_PROTOCOL_UNSPECIFIED)
	}

	return out
}

func PortBindingFromSpec(spec string) (PortBinding, error) {
	var strHostPort, strContainerPort, strIP, strProtocol string
	var intHostPort, intContainerPort uint64
	var err error

	out := PortBinding{}

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return out, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	case 1:
		strContainerPort = values[0]
	case 2:
		strHostPort = values[0]
		strContainerPort = values[1]
	case 3:
		strIP = values[0]
		strHostPort = values[1]
		strContainerPort = values[2]
	default:
		return out, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	switch values := strings.SplitN(strHostPort, "/", 2); len(values) {
	case 1:
		strHostPort = values[0]
	case 2:
		strHostPort = values[0]
		strProtocol = values[1]
	default:
		return out, fmt.Errorf("%s: %w", spec, ErrPortFormat)
	}

	if strContainerPort != "" {
		intContainerPort, err = strconv.ParseUint(strContainerPort, 10, 16)
		if err != nil {
			return out, fmt.Errorf("%s: %w", spec, ErrPortFormat)
		}

		out.ContainerPort = uint16(intContainerPort)
	}

	if strHostPort != "" {
		intHostPort, err = strconv.ParseUint(strHostPort, 10, 16)
		if err != nil {
			return out, fmt.Errorf("%s: %w", spec, ErrPortFormat)
		}

		out.HostPort = uint16(intHostPort)
	}

	if strProtocol != "" {
		switch strProtocol {
		case "tcp":
			out.Protocol = TCP
		case "udp":
			out.Protocol = UDP
		default:
			return out, fmt.Errorf("%s: %w", spec, ErrPortFormat)
		}
	}

	if strIP != "" {
		out.HostIP = strIP
	}

	return out, nil
}
