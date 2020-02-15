package configuration

import (
	"github.com/oclaussen/dodo/pkg/plugin/configuration/proto"
	"github.com/oclaussen/dodo/pkg/types"
	"golang.org/x/net/context"
)

type client struct {
	configClient proto.ConfigurationClient
}

func (c *client) GetClientOptions(backdrop string) (*ClientOptions, error) {
	opts, err := c.configClient.GetClientOptions(context.Background(), &proto.Backdrop{Name: backdrop})
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
	return c.configClient.UpdateConfiguration(context.Background(), backdrop)
}

func (c *client) Provision(containerID string) error {
	_, err := c.configClient.Provision(context.Background(), &proto.Container{Id: containerID})
	return err
}
