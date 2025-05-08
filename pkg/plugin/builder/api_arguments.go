package builder

import (
	"errors"
	"fmt"
	"os"
	"strings"

	api "github.com/wabenet/dodo-core/api/build/v1alpha2"
)

var ErrArgumentFormat = errors.New("invalid argument format")

type BuildArgumentsConfig []BuildArgument

type BuildArgument struct {
	Key   string
	Value string
}

func MergeBuildArgumentsConfig(first, second BuildArgumentsConfig) BuildArgumentsConfig {
	return append(first, second...)
}

func BuildArgumentsConfigFromProto(b []*api.BuildArgument) BuildArgumentsConfig {
	out := BuildArgumentsConfig{}

	for _, arg := range b {
		out = append(out, BuildArgumentFromProto(arg))
	}

	return out
}

func (b BuildArgumentsConfig) ToProto() []*api.BuildArgument {
	out := []*api.BuildArgument{}

	for _, arg := range b {
		out = append(out, arg.ToProto())
	}

	return out
}

func BuildArgumentFromProto(b *api.BuildArgument) BuildArgument {
	return BuildArgument{
		Key:   b.GetKey(),
		Value: b.GetValue(),
	}
}

func (b BuildArgument) ToProto() *api.BuildArgument {
	out := &api.BuildArgument{}

	out.SetKey(b.Key)
	out.SetValue(b.Value)

	return out
}

func BuildArgumentFromSpec(spec string) (BuildArgument, error) {
	out := BuildArgument{}

	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 0:
		return out, fmt.Errorf("%s: %w", spec, ErrArgumentFormat)
	case 1:
		out.Key = values[0]
		out.Value = os.Getenv(values[0])
	case 2:
		out.Key = values[0]
		out.Value = values[1]
	default:
		return out, fmt.Errorf("%s: %w", spec, ErrArgumentFormat)
	}

	return out, nil
}
