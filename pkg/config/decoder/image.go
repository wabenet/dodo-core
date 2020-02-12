package decoder

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/oclaussen/dodo/pkg/types"
)

func (d *decoder) DecodeImage(name string, config interface{}) (types.BuildInfo, error) {
	var result types.BuildInfo
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		decoded, err := d.DecodeString(name, config)
		if err != nil {
			return result, err
		}
		result.ImageName = decoded
	case reflect.Map:
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			switch key := k.(string); key {
			case "name":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return result, err
				}
				result.ImageName = decoded
			case "context":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return result, err
				}
				result.Context = decoded
			case "dockerfile":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return result, err
				}
				result.Dockerfile = decoded
			case "steps", "inline":
				decoded, err := d.DecodeStringSlice(key, v)
				if err != nil {
					return result, err
				}
				result.InlineDockerfile = decoded
			case "args", "arguments":
				decoded, err := d.DecodeArguments(key, v)
				if err != nil {
					return result, err
				}
				result.Arguments = decoded
			case "secrets":
				decoded, err := d.DecodeSecrets(key, v)
				if err != nil {
					return result, err
				}
				result.Secrets = decoded
			case "ssh":
				decoded, err := d.DecodeSSHAgents(key, v)
				if err != nil {
					return result, err
				}
				result.SshAgents = decoded
			case "no_cache":
				decoded, err := d.DecodeBool(key, v)
				if err != nil {
					return result, err
				}
				result.NoCache = decoded
			case "force_rebuild":
				decoded, err := d.DecodeBool(key, v)
				if err != nil {
					return result, err
				}
				result.ForceRebuild = decoded
			case "force_pull":
				decoded, err := d.DecodeBool(key, v)
				if err != nil {
					return result, err
				}
				result.ForcePull = decoded
			case "requires", "dependencies":
				decoded, err := d.DecodeStringSlice(key, v)
				if err != nil {
					return result, err
				}
				result.Dependencies = decoded
			}
		}
	default:
		return result, &ConfigError{Name: name, UnsupportedType: t.Kind()}
	}
	return result, nil
}

func (d *decoder) DecodeArguments(name string, config interface{}) ([]*types.Argument, error) {
	result := []*types.Argument{}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		decoded, err := d.DecodeArgument(name, config)
		if err != nil {
			return result, err
		}
		result = append(result, decoded)
	case reflect.Slice:
		for _, v := range t.Interface().([]interface{}) {
			decoded, err := d.DecodeArguments(name, v)
			if err != nil {
				return result, err
			}
			result = append(result, decoded...)
		}
	case reflect.Map:
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			key := k.(string)
			decoded, err := d.DecodeString(key, v)
			if err != nil {
				return result, err
			}
			result = append(result, &types.Argument{
				Key:   key,
				Value: decoded,
			})
		}
	default:
		return result, &ConfigError{Name: name, UnsupportedType: t.Kind()}
	}
	return result, nil
}

func (d *decoder) DecodeArgument(name string, config interface{}) (*types.Argument, error) {
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		decoded, err := d.DecodeString(name, t.String())
		if err != nil {
			return nil, err
		}
		switch values := strings.SplitN(decoded, "=", 2); len(values) {
		case 0:
			return nil, fmt.Errorf("empty assignment in '%s'", name)
		case 1:
			return &types.Argument{Key: values[0], Value: os.Getenv(values[0])}, nil
		case 2:
			return &types.Argument{Key: values[0], Value: values[1]}, nil
		default:
			return nil, fmt.Errorf("too many values in '%s'", name)
		}
	default:
		return nil, &ConfigError{Name: name, UnsupportedType: t.Kind()}
	}
}
