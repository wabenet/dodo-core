package types

import (
	"reflect"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/appconfig"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

func NewDevice() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.DeviceMapping{}
		return &target, DecodeDevice(&target)
	}
}

func DecodeDevice(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	dev := *(target.(**api.DeviceMapping))

	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"cgroup_rule": decoder.String(&dev.CgroupRule),
			"source":      decoder.String(&dev.Source),
			"target":      decoder.String(&dev.Target),
			"permissions": decoder.String(&dev.Permissions),
		}),
		reflect.String: func(d *decoder.Decoder, config interface{}) {
			var decoded string
			decoder.String(&decoded)(d, config)

			if dv, err := appconfig.ParseDeviceMapping(decoded); err != nil {
				d.Error(err)
			} else {
				dev.CgroupRule = dv.CgroupRule
				dev.Source = dv.Source
				dev.Target = dv.Target
				dev.Permissions = dv.Permissions
			}
		},
	})
}
