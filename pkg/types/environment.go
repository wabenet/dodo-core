package types

import (
	"fmt"
	"os"
	"strings"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

const ErrEnvironmentFormat FormatError = "invalid environment format"

func ParseEnvironment(spec string) (*api.Environment, error) {
	env := &api.Environment{}

	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 0:
		return nil, fmt.Errorf("%s: %w", spec, ErrEnvironmentFormat)
	case 1:
		env.Key = values[0]
		env.Value = os.Getenv(values[0])
	case 2:
		env.Key = values[0]
		env.Value = values[1]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrEnvironmentFormat)
	}

	return env, nil
}

func NewEnvironment() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.Environment{}
		return &target, DecodeEnvironment(&target)
	}
}

func DecodeEnvironment(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	env := *(target.(**api.Environment))

	return func(d *decoder.Decoder, config interface{}) {
		var decoded string

		decoder.String(&decoded)(d, config)

		if e, err := ParseEnvironment(decoded); err != nil {
			d.Error(err)
		} else {
			env.Key = e.Key
			env.Value = e.Value
		}
	}
}
