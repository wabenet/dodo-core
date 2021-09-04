package types

import (
	"fmt"
	"reflect"
	"strings"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

const ErrPortFormat FormatError = "invalid publish format"

func ParsePort(spec string) (*api.Port, error) {
	port := &api.Port{}

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

	switch values := strings.SplitN(port.Target, "/", 2); len(values) {
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

func NewPort() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.Port{}
		return &target, DecodePort(&target)
	}
}

func DecodePort(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	port := *(target.(**api.Port))

	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"target":    decoder.String(&port.Target),
			"published": decoder.String(&port.Published),
			"protocol":  decoder.String(&port.Protocol),
			"host_ip":   decoder.String(&port.HostIp),
		}),
		reflect.String: func(d *decoder.Decoder, config interface{}) {
			var decoded string
			decoder.String(&decoded)(d, config)
			if p, err := ParsePort(decoded); err != nil {
				d.Error(err)
			} else {
				port.Target = p.Target
				port.Published = p.Published
				port.Protocol = p.Protocol
				port.HostIp = p.HostIp
			}
		},
	})
}
