package decoder

import (
	"reflect"

	"github.com/oclaussen/dodo/pkg/types"
)

func (d *decoder) DecodeSSHAgents(name string, config interface{}) ([]*types.SshAgent, error) {
	result := []*types.SshAgent{}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.Bool:
		decoded, err := d.DecodeBool(name, config)
		if err != nil {
			return result, err
		}
		if decoded {
			result = append(result, &types.SshAgent{})
		}
	case reflect.String, reflect.Map:
		decoded, err := d.DecodeSSHAgent(name, config)
		if err != nil {
			return result, err
		}
		result = append(result, decoded)
	case reflect.Slice:
		for _, v := range t.Interface().([]interface{}) {
			decoded, err := d.DecodeSSHAgent(name, v)
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

func (d *decoder) DecodeSSHAgent(name string, config interface{}) (*types.SshAgent, error) {
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		decoded, err := d.DecodeEnvironment(name, config) // FIXME: lazy hack
		if err != nil {
			return nil, err
		}
		return &types.SshAgent{Id: decoded.Key, IdentityFile: decoded.Value}, nil
	case reflect.Map:
		var result types.SshAgent
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			switch key := k.(string); key {
			case "id":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.Id = decoded
			case "file":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.IdentityFile = decoded
			}
		}
		return &result, nil
	default:
		return nil, &ConfigError{Name: name, UnsupportedType: t.Kind()}
	}
}
