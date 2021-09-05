package types

import (
	"reflect"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/appconfig"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

func NewVolume() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.VolumeMount{}
		return &target, DecodeVolume(&target)
	}
}

func DecodeVolume(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	vol := *(target.(**api.VolumeMount))

	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"source":    decoder.String(&vol.Source),
			"target":    decoder.String(&vol.Target),
			"read_only": decoder.Bool(&vol.Readonly),
		}),
		reflect.String: func(d *decoder.Decoder, config interface{}) {
			var decoded string
			decoder.String(&decoded)(d, config)

			if v, err := appconfig.ParseVolumeMount(decoded); err != nil {
				d.Error(err)
			} else {
				vol.Source = v.Source
				vol.Target = v.Target
				vol.Readonly = v.Readonly
			}
		},
	})
}
