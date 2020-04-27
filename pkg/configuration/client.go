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
	backdrop.Merge(update)
	return backdrop, nil
}

func (c *client) Provision(containerID string) error {
	_, err := c.configClient.Provision(context.Background(), &types.ContainerId{Id: containerID})
	return err
}
