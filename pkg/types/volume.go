package types

import (
	"fmt"
	"reflect"
	"strings"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

const ErrVolumeFormat FormatError = "invalid volume format"

func ParseVolume(spec string) (*api.Volume, error) {
	vol := &api.Volume{}

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

func NewVolume() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.Volume{}
		return &target, DecodeVolume(&target)
	}
}

func DecodeVolume(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	vol := *(target.(**api.Volume))

	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"source":    decoder.String(&vol.Source),
			"target":    decoder.String(&vol.Target),
			"read_only": decoder.Bool(&vol.Readonly),
		}),
		reflect.String: func(d *decoder.Decoder, config interface{}) {
			var decoded string
			decoder.String(&decoded)(d, config)

			if v, err := ParseVolume(decoded); err != nil {
				d.Error(err)
			} else {
				vol.Source = v.Source
				vol.Target = v.Target
				vol.Readonly = v.Readonly
			}
		},
	})
}
