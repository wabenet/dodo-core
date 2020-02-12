package decoder

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/oclaussen/dodo/pkg/types"
)

func (d *decoder) DecodeVolumes(name string, config interface{}) ([]*types.Volume, error) {
	result := []*types.Volume{}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String, reflect.Map:
		decoded, err := d.DecodeVolume(name, config)
		if err != nil {
			return result, err
		}
		result = append(result, decoded)
	case reflect.Slice:
		for _, v := range t.Interface().([]interface{}) {
			decoded, err := d.DecodeVolume(name, v)
			if err != nil {
				return result, err
			}
			result = append(result, decoded)
		}
	default:
		return result, &ConfigError{Name: name, UnsupportedType: t.Kind()}
	}
	return result, nil
}

func (d *decoder) DecodeVolume(name string, config interface{}) (*types.Volume, error) {
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		decoded, err := d.DecodeString(name, t.String())
		if err != nil {
			return nil, err
		}
		switch values := strings.SplitN(decoded, ":", 3); len(values) {
		case 0:
			return nil, fmt.Errorf("empty volume definition in '%s'", name)
		case 1:
			return &types.Volume{
				Source: values[0],
			}, nil
		case 2:
			return &types.Volume{
				Source: values[0],
				Target: values[1],
			}, nil
		case 3:
			return &types.Volume{
				Source:   values[0],
				Target:   values[1],
				Readonly: values[2] == "ro",
			}, nil
		default:
			return nil, fmt.Errorf("too many values in '%s'", name)
		}
	case reflect.Map:
		var result types.Volume
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			switch key := k.(string); key {
			case "source":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.Source = decoded
			case "target":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.Target = decoded
			case "read_only":
				decoded, err := d.DecodeBool(key, v)
				if err != nil {
					return nil, err
				}
				result.Readonly = decoded
			}
		}
		return &result, nil
	default:
		return nil, &ConfigError{Name: name, UnsupportedType: t.Kind()}
	}
}
