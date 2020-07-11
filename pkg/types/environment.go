package types

import (
	"fmt"
	"os"
	"strings"

	"github.com/oclaussen/dodo/pkg/decoder"
)

func (env *Environment) FromString(spec string) error {
	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 0:
		return fmt.Errorf("empty assignment in environment: %s", spec)
	case 1:
		env.Key = values[0]
		env.Value = os.Getenv(values[0])
	case 2:
		env.Key = values[0]
		env.Value = values[1]
	default:
		return fmt.Errorf("too many values in environment definition: %s", spec)
	}

	return nil
}

func NewEnvironment() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &Environment{}
		return &target, DecodeEnvironment(&target)
	}
}

func DecodeEnvironment(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	env := *(target.(**Environment))
	return func(d *decoder.Decoder, config interface{}) {
		var decoded string
		decoder.String(&decoded)(d, config)
		if err := env.FromString(decoded); err != nil {
			d.Error("invalid environment")
		}
	}
}
