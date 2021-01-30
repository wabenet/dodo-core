package types

import (
	"reflect"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

func NewBackdrop() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.Backdrop{Entrypoint: &api.Entrypoint{}}
		return &target, DecodeBackdrop(&target)
	}
}

func DecodeBackdrop(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	backdrop := *(target.(**api.Backdrop))

	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"name":           decoder.String(&backdrop.ContainerName),
			"alias":          decoder.Slice(decoder.NewString(), &backdrop.Aliases),
			"aliases":        decoder.Slice(decoder.NewString(), &backdrop.Aliases),
			"container_name": decoder.String(&backdrop.ContainerName),
			"image":          decoder.String(&backdrop.ImageId),
			"runtime":        decoder.String(&backdrop.Runtime),
			"interactive":    decoder.Bool(&backdrop.Entrypoint.Interactive),
			"script":         decoder.String(&backdrop.Entrypoint.Script),
			"interpreter": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(decoder.NewString(), &backdrop.Entrypoint.Interpreter),
				reflect.Slice:  decoder.Slice(decoder.NewString(), &backdrop.Entrypoint.Interpreter),
			}),
			"arguments": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(decoder.NewString(), &backdrop.Entrypoint.Arguments),
				reflect.Slice:  decoder.Slice(decoder.NewString(), &backdrop.Entrypoint.Arguments),
			}),
			"command": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(decoder.NewString(), &backdrop.Entrypoint.Arguments),
				reflect.Slice:  decoder.Slice(decoder.NewString(), &backdrop.Entrypoint.Arguments),
			}),
			"env": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewEnvironment(), &backdrop.Environment),
				reflect.Slice:  decoder.Slice(NewEnvironment(), &backdrop.Environment),
			}),
			"environment": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewEnvironment(), &backdrop.Environment),
				reflect.Slice:  decoder.Slice(NewEnvironment(), &backdrop.Environment),
			}),
			"volume": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewVolume(), &backdrop.Volumes),
				reflect.Map:    decoder.Singleton(NewVolume(), &backdrop.Volumes),
				reflect.Slice:  decoder.Slice(NewVolume(), &backdrop.Volumes),
			}),
			"volumes": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewVolume(), &backdrop.Volumes),
				reflect.Map:    decoder.Singleton(NewVolume(), &backdrop.Volumes),
				reflect.Slice:  decoder.Slice(NewVolume(), &backdrop.Volumes),
			}),
			"device": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewDevice(), &backdrop.Devices),
				reflect.Map:    decoder.Singleton(NewDevice(), &backdrop.Devices),
				reflect.Slice:  decoder.Slice(NewDevice(), &backdrop.Devices),
			}),
			"devices": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewDevice(), &backdrop.Devices),
				reflect.Map:    decoder.Singleton(NewDevice(), &backdrop.Devices),
				reflect.Slice:  decoder.Slice(NewDevice(), &backdrop.Devices),
			}),
			"ports": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewPort(), &backdrop.Ports),
				reflect.Map:    decoder.Singleton(NewPort(), &backdrop.Ports),
				reflect.Slice:  decoder.Slice(NewPort(), &backdrop.Ports),
			}),
			"capabilities": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(decoder.NewString(), &backdrop.Capabilities),
				reflect.Slice:  decoder.Slice(decoder.NewString(), &backdrop.Capabilities),
			}),
			"user":        decoder.String(&backdrop.User),
			"workdir":     decoder.String(&backdrop.WorkingDir),
			"working_dir": decoder.String(&backdrop.WorkingDir),
		}),
	})
}

func Merge(target *api.Backdrop, source *api.Backdrop) {
	if len(source.Name) > 0 {
		target.Name = source.Name
	}

	target.Aliases = append(target.Aliases, source.Aliases...)

	if len(source.ImageId) > 0 {
		target.ImageId = source.ImageId
	}

	if source.Entrypoint != nil {
		if source.Entrypoint.Interactive {
			target.Entrypoint.Interactive = true
		}

		if len(source.Entrypoint.Interpreter) > 0 {
			target.Entrypoint.Interpreter = source.Entrypoint.Interpreter
		}

		if len(source.Entrypoint.Script) > 0 {
			target.Entrypoint.Script = source.Entrypoint.Script
		}

		if len(source.Entrypoint.Arguments) > 0 {
			target.Entrypoint.Arguments = source.Entrypoint.Arguments
		}
	}

	if len(source.ContainerName) > 0 {
		target.ContainerName = source.ContainerName
	}

	target.Environment = append(target.Environment, source.Environment...)

	if len(source.User) > 0 {
		target.User = source.User
	}

	target.Volumes = append(target.Volumes, source.Volumes...)
	target.Devices = append(target.Devices, source.Devices...)
	target.Ports = append(target.Ports, source.Ports...)
	target.Capabilities = append(target.Capabilities, source.Capabilities...)

	if len(source.WorkingDir) > 0 {
		target.WorkingDir = source.WorkingDir
	}
}
