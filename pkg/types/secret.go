package types

import (
	"reflect"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

func NewSecret() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.Secret{}
		return &target, DecodeSecret(&target)
	}
}

func DecodeSecret(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	secret := *(target.(**api.Secret))

	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"id":     decoder.String(&secret.Id),
			"src":    decoder.String(&secret.Path),
			"source": decoder.String(&secret.Path),
		}),
	})
}
