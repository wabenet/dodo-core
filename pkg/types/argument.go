package types

import (
	"fmt"
	"os"
	"strings"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

const ErrArgumentFormat FormatError = "invalid argument format"

func ParseArgument(spec string) (*api.Argument, error) {
	arg := &api.Argument{}

	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 0:
		return nil, fmt.Errorf("%s: %w", spec, ErrArgumentFormat)
	case 1:
		arg.Key = values[0]
		arg.Value = os.Getenv(values[0])
	case 2:
		arg.Key = values[0]
		arg.Value = values[1]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrArgumentFormat)
	}

	return arg, nil
}

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

		if a, err := ParseArgument(decoded); err != nil {
			d.Error(err)
		} else {
			arg.Key = a.Key
			arg.Value = a.Value
		}
	}
}
