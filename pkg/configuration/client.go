package configuration

import (
	"github.com/oclaussen/dodo/pkg/types"
	"golang.org/x/net/context"
)

type client struct {
	configClient types.ConfigurationClient
}

func (c *client) Init() error {
	_, err := c.configClient.Init(context.Background(), &types.Empty{})
	return err
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
