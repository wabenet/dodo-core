package types

import (
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/config"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

func NewEnvironment() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.EnvironmentVariable{}
		return &target, DecodeEnvironment(&target)
	}
}

func DecodeEnvironment(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	env := *(target.(**api.EnvironmentVariable))

	return func(d *decoder.Decoder, c interface{}) {
		var decoded string

		decoder.String(&decoded)(d, c)

		if e, err := config.ParseEnvironmentVariable(decoded); err != nil {
			d.Error(err)
		} else {
			env.Key = e.Key
			env.Value = e.Value
		}
	}
}
