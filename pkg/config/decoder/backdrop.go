package decoder

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/oclaussen/dodo/pkg/types"
)

func (d *decoder) DecodeBackdrops(name string, config interface{}) (map[string]types.Backdrop, error) {
	result := map[string]types.Backdrop{}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.Map:
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			key := k.(string)
			decoded, err := d.DecodeBackdrop(key, v)
			if err != nil {
				return result, err
			}
			result[key] = decoded
			for _, alias := range decoded.Aliases {
				result[alias] = decoded
			}
		}
	default:
		return result, &ConfigError{Name: name, UnsupportedType: t.Kind()}
	}
	return result, nil
}

func (d *decoder) DecodeBackdrop(name string, config interface{}) (types.Backdrop, error) {
	result := types.Backdrop{Name: name, Entrypoint: &types.Entrypoint{}} // TODO , filename: d.filename}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.Map:
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			switch key := k.(string); key {
			case "alias", "aliases":
				decoded, err := d.DecodeStringSlice(key, v)
				if err != nil {
					return result, err
				}
				result.Aliases = decoded
			case "build", "image":
				decoded, err := d.DecodeImage(key, v)
				if err != nil {
					return result, err
				}
				result.Build = &decoded
			case "name", "container_name":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return result, err
				}
				result.ContainerName = decoded
			case "interactive":
				decoded, err := d.DecodeBool(key, v)
				if err != nil {
					return result, err
				}
				result.Entrypoint.Interactive = decoded
			case "env", "environment":
				decoded, err := d.DecodeEnvironments(key, v)
				if err != nil {
					return result, err
				}
				result.Environment = decoded
			case "user":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return result, err
				}
				result.User = decoded
			case "volume", "volumes":
				decoded, err := d.DecodeVolumes(key, v)
				if err != nil {
					return result, err
				}
				result.Volumes = decoded
			case "device", "devices":
				decoded, err := d.DecodeDevices(key, v)
				if err != nil {
					return result, err
				}
				result.Devices = decoded
			case "ports":
				decoded, err := d.DecodePorts(key, v)
				if err != nil {
					return result, err
				}
				result.Ports = decoded
			case "workdir", "working_dir":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return result, err
				}
				result.WorkingDir = decoded
			case "interpreter":
				decoded, err := d.DecodeStringSlice(key, v)
				if err != nil {
					return result, err
				}
				result.Entrypoint.Interpreter = decoded
			case "script":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return result, err
				}
				result.Entrypoint.Script = decoded
			case "command", "arguments":
				decoded, err := d.DecodeStringSlice(key, v)
				if err != nil {
					return result, err
				}
				result.Entrypoint.Arguments = decoded
			default:
				return result, &ConfigError{Name: name, UnsupportedKey: &key}
			}
		}
	default:
		return result, &ConfigError{Name: name, UnsupportedType: t.Kind()}
	}
	return result, nil
}

func (d *decoder) DecodeEnvironments(name string, config interface{}) ([]*types.Environment, error) {
	result := []*types.Environment{}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		decoded, err := d.DecodeEnvironment(name, config)
		if err != nil {
			return result, err
		}
		result = append(result, decoded)
	case reflect.Slice:
		for _, v := range t.Interface().([]interface{}) {
			decoded, err := d.DecodeEnvironments(name, v)
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
			result = append(result, &types.Environment{
				Key:   key,
				Value: decoded,
			})
		}
	default:
		return result, &ConfigError{Name: name, UnsupportedType: t.Kind()}
	}
	return result, nil
}

func (d *decoder) DecodeEnvironment(name string, config interface{}) (*types.Environment, error) {
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
			return &types.Environment{Key: values[0], Value: os.Getenv(values[0])}, nil
		case 2:
			return &types.Environment{Key: values[0], Value: values[1]}, nil
		default:
			return nil, fmt.Errorf("too many values in '%s'", name)
		}
	default:
		return nil, &ConfigError{Name: name, UnsupportedType: t.Kind()}
	}
}
