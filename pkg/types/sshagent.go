package types

import (
	"reflect"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/config"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

func NewSSHAgent() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.SshAgent{}
		return &target, DecodeSSHAgent(&target)
	}
}

func DecodeSSHAgent(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	agent := *(target.(**api.SshAgent))

	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"id":   decoder.String(&agent.Id),
			"file": decoder.String(&agent.IdentityFile),
		}),
		reflect.String: func(d *decoder.Decoder, c interface{}) {
			var decoded string

			decoder.String(&decoded)(d, c)

			if a, err := config.ParseSSHAgent(decoded); err != nil {
				d.Error(err)
			} else {
				agent.Id = a.Id
				agent.IdentityFile = a.IdentityFile
			}
		},
	})
}
