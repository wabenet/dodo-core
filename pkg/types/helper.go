package types

import (
	"fmt"
	"os"
	"strings"
)

func NewVolume(spec string) (*Volume, error) {
	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return nil, fmt.Errorf("empty volume definition: %s", spec)
	case 1:
		return &Volume{Source: values[0]}, nil
	case 2:
		return &Volume{
			Source: values[0],
			Target: values[1],
		}, nil
	case 3:
		return &Volume{
			Source:   values[0],
			Target:   values[1],
			Readonly: values[2] == "ro",
		}, nil
	default:
		return nil, fmt.Errorf("too many values in volume definition: %s", spec)
	}
}

func NewEnvironment(spec string) (*Environment, error) {
	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 0:
		return nil, fmt.Errorf("empty assignment in environment: %s", spec)
	case 1:
		return &Environment{
			Key:   values[0],
			Value: os.Getenv(values[0]),
		}, nil
	case 2:
		return &Environment{
			Key:   values[0],
			Value: values[1],
		}, nil
	default:
		return nil, fmt.Errorf("too many values in environment definition: %s", spec)
	}
}

func NewDevice(spec string) (*Device, error) {
	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return nil, fmt.Errorf("empty device definition: %s", spec)
	case 1:
		return &Device{
			Source: values[0],
		}, nil
	case 2:
		return &Device{
			Source: values[0],
			Target: values[1],
		}, nil
	case 3:
		return &Device{
			Source:      values[0],
			Target:      values[1],
			Permissions: values[2],
		}, nil
	default:
		return nil, fmt.Errorf("too many values in device definition: %s", spec)
	}
}

func NewPort(spec string) (*Port, error) {
	var result *Port

	switch values := strings.SplitN(spec, ":", 3); len(values) {
	case 0:
		return result, fmt.Errorf("empty publish definition: %s", spec)
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
		return nil, fmt.Errorf("too many values in publish definition: %s", spec)
	}

	switch values := strings.SplitN(result.Target, "/", 2); len(values) {
	case 1:
		result.Target = values[0]
	case 2:
		result.Target = values[0]
		result.Protocol = values[1]
	default:
		return nil, fmt.Errorf("too many values in publish definition: %s", spec)
	}

	return result, nil
}
