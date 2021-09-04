package types

import (
	"fmt"
	"reflect"
	"strings"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

const ErrDeviceFormat FormatError = "invalid device format"

func ParseDevice(spec string) (*api.Device, error) {
	dev := &api.Device{}

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

func NewDevice() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.Device{}
		return &target, DecodeDevice(&target)
	}
}

func DecodeDevice(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	dev := *(target.(**api.Device))

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

			if dv, err := ParseDevice(decoded); err != nil {
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
