package runtime

import (
	"errors"
	"fmt"
	"os"
	"strings"

	api "github.com/wabenet/dodo-core/api/runtime/v1alpha2"
)

var ErrEnvironmentFormat = errors.New("invalid environment format")

type EnvironmentConfig []EnvironmentVariable

type EnvironmentVariable struct {
	Key   string
	Value string
}

func MergeEnvironmentConfig(first, second EnvironmentConfig) EnvironmentConfig {
	return append(first, second...)
}

func (e EnvironmentVariable) String() string {
	return fmt.Sprintf("%s=%s", e.Key, e.Value)
}

func EnvironmentConfigFromProto(e []*api.EnvironmentVariable) EnvironmentConfig {
	out := EnvironmentConfig{}

	for _, env := range e {
		out = append(out, EnvironmentVariableFromProto(env))
	}

	return out
}

func EnvironmentVariableFromProto(e *api.EnvironmentVariable) EnvironmentVariable {
	return EnvironmentVariable{
		Key:   e.GetKey(),
		Value: e.GetValue(),
	}
}

func (e EnvironmentConfig) ToProto() []*api.EnvironmentVariable {
	out := []*api.EnvironmentVariable{}

	for _, env := range e {
		out = append(out, env.ToProto())
	}

	return out
}

func (e EnvironmentVariable) ToProto() *api.EnvironmentVariable {
	out := &api.EnvironmentVariable{}

	out.SetKey(e.Key)
	out.SetValue(e.Value)

	return out
}

func EnvironmentVariableFromSpec(spec string) (EnvironmentVariable, error) {
	out := EnvironmentVariable{}

	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 0:
		return out, fmt.Errorf("%s: %w", spec, ErrEnvironmentFormat)
	case 1:
		out.Key = values[0]
		out.Value = os.Getenv(values[0])
	case 2:
		out.Key = values[0]
		out.Value = values[1]
	default:
		return out, fmt.Errorf("%s: %w", spec, ErrEnvironmentFormat)
	}

	return out, nil
}
