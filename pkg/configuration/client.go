package configuration

import (
	"github.com/oclaussen/dodo/pkg/types"
	"golang.org/x/net/context"
)

type client struct {
	configClient types.ConfigurationClient
}

func (c *client) GetClientOptions(backdrop string) (*ClientOptions, error) {
	opts, err := c.configClient.GetClientOptions(context.Background(), &types.Backdrop{Name: backdrop})
	if err != nil {
		return nil, err
	}
	return &ClientOptions{
		Version:  opts.Version,
		Host:     opts.Host,
		CAFile:   opts.CaFile,
		CertFile: opts.CertFile,
		KeyFile:  opts.KeyFile,
	}, nil
}

func (c *client) UpdateConfiguration(backdrop *types.Backdrop) (*types.Backdrop, error) {
	update, err := c.configClient.UpdateConfiguration(context.Background(), backdrop)
	if err != nil {
		return nil, err
	}
	mergeBackdrops(backdrop, update)
	return backdrop, nil
}

func (c *client) Provision(containerID string) error {
	_, err := c.configClient.Provision(context.Background(), &types.ContainerId{Id: containerID})
	return err
}

func mergeBackdrops(target *types.Backdrop, source *types.Backdrop) {
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
	if len(source.WorkingDir) > 0 {
		target.WorkingDir = source.WorkingDir
	}
}
