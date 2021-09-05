package types

import (
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/appconfig"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

func NewArgument() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.Argument{}
		return &target, DecodeArgument(&target)
	}
}

func DecodeArgument(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	arg := *(target.(**api.Argument))

	return func(d *decoder.Decoder, config interface{}) {
		var decoded string

		decoder.String(&decoded)(d, config)

		if a, err := appconfig.ParseArgument(decoded); err != nil {
			d.Error(err)
		} else {
			arg.Key = a.Key
			arg.Value = a.Value
		}
	}
}
