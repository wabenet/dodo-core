package decoder

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/oclaussen/dodo/pkg/types"
)

func (d *decoder) DecodePorts(name string, config interface{}) ([]*types.Port, error) {
	result := []*types.Port{}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String, reflect.Map:
		decoded, err := d.DecodePort(name, config)
		if err != nil {
			return result, err
		}
		result = append(result, decoded)
	case reflect.Slice:
		for _, v := range t.Interface().([]interface{}) {
			decoded, err := d.DecodePort(name, v)
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

func (d *decoder) DecodePort(name string, config interface{}) (*types.Port, error) {
	result := types.Port{Protocol: "tcp"}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		decoded, err := d.DecodeString(name, t.String())
		if err != nil {
			return nil, err
		}
		switch values := strings.SplitN(decoded, ":", 3); len(values) {
		case 0:
			return nil, fmt.Errorf("empty port definition in '%s'", name)
		case 1:
			result.Target = values[0]
		case 2:
			result.Published = values[0]
			result.Target = values[1]
		case 3:
			result.HostIp = values[0]
			result.Published = values[1]
			result.Target = values[2]
		default:
			return nil, fmt.Errorf("too many values in '%s'", name)
		}
		switch values := strings.SplitN(result.Target, "/", 2); len(values) {
		case 1:
			result.Target = values[0]
		case 2:
			result.Target = values[0]
			result.Protocol = values[1]
		default:
			return nil, fmt.Errorf("too many values in '%s'", name)
		}
		return &result, nil
	case reflect.Map:
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			switch key := k.(string); key {
			case "target":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.Target = decoded
			case "published":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.Published = decoded
			case "protocol":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.Protocol = decoded
			case "host_ip":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.HostIp = decoded
			default:
				return nil, &ConfigError{Name: name, UnsupportedKey: &key}
			}
		}
		return &result, nil
	default:
		return nil, &ConfigError{Name: name, UnsupportedType: t.Kind()}
	}
}
