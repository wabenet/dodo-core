package types

import (
	"reflect"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/appconfig"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

func NewPort() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.PortBinding{}
		return &target, DecodePort(&target)
	}
}

func DecodePort(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	port := *(target.(**api.PortBinding))

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
			if p, err := appconfig.ParsePortBinding(decoded); err != nil {
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
