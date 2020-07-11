package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/oclaussen/dodo/pkg/decoder"
)

func (vol *Volume) FromString(spec string) error {
	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return fmt.Errorf("empty volume definition: %s", spec)
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
		return fmt.Errorf("too many values in volume definition: %s", spec)
	}

	return nil
}

func NewVolume() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &Volume{}
		return &target, DecodeVolume(&target)
	}
}

func DecodeVolume(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	vol := *(target.(**Volume))
	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"source":    decoder.String(&vol.Source),
			"target":    decoder.String(&vol.Target),
			"read_only": decoder.Bool(&vol.Readonly),
		}),
		reflect.String: func(d *decoder.Decoder, config interface{}) {
			var decoded string
			decoder.String(&decoded)(d, config)
			if err := vol.FromString(decoded); err != nil {
				d.Error("invalid volume")
			}
		},
	})
}
