package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/oclaussen/dodo/pkg/decoder"
)

func (dev *Device) FromString(spec string) error {
	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return fmt.Errorf("empty device definition: %s", spec)
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
		return fmt.Errorf("too many values in device definition: %s", spec)
	}

	return nil
}

func NewDevice() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &Device{}
		return &target, DecodeDevice(&target)
	}
}

func DecodeDevice(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	dev := *(target.(**Device))
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
			if err := dev.FromString(decoded); err != nil {
				d.Error("invalid device")
			}
		},
	})
}
