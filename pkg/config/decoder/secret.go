package decoder

import (
	"encoding/csv"
	"reflect"
	"strings"

	"github.com/oclaussen/dodo/pkg/types"
)

func (d *decoder) DecodeSecrets(name string, config interface{}) ([]*types.Secret, error) {
	result := []*types.Secret{}
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String, reflect.Map:
		decoded, err := d.DecodeSecret(name, config)
		if err != nil {
			return result, err
		}
		result = append(result, decoded)
	case reflect.Slice:
		for _, v := range t.Interface().([]interface{}) {
			decoded, err := d.DecodeSecret(name, v)
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

func (d *decoder) DecodeSecret(name string, config interface{}) (*types.Secret, error) {
	switch t := reflect.ValueOf(config); t.Kind() {
	case reflect.String:
		decoded, err := d.DecodeString(name, t.String())
		if err != nil {
			return nil, err
		}

		reader := csv.NewReader(strings.NewReader(decoded))
		fields, err := reader.Read()
		if err != nil {
			return nil, err
		}

		secretMap := make(map[interface{}]interface{}, len(fields))
		for _, field := range fields {
			kv, err := d.DecodeEnvironment(name, field) // FIXME: lazy hack
			if err != nil {
				return nil, err
			}
			secretMap[kv.Key] = kv.Value
		}
		return d.DecodeSecret(name, secretMap)
	case reflect.Map:
		var result types.Secret
		for k, v := range t.Interface().(map[interface{}]interface{}) {
			switch key := k.(string); key {
			case "id":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.Id = decoded
			case "source", "src":
				decoded, err := d.DecodeString(key, v)
				if err != nil {
					return nil, err
				}
				result.Path = decoded
			}
		}
		return &result, nil
	default:
		return nil, &ConfigError{Name: name, UnsupportedType: t.Kind()}
	}
}
