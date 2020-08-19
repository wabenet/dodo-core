package configuration

import (
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/types"
	"golang.org/x/net/context"
)

type client struct {
	configClient types.ConfigurationClient
}

func (t *client) Type() plugin.Type {
	return Type
}

func (c *client) Init() error {
	_, err := c.configClient.Init(context.Background(), &types.Empty{})
	return err
}

func (c *client) UpdateConfiguration(backdrop *types.Backdrop) (*types.Backdrop, error) {
	return c.configClient.UpdateConfiguration(context.Background(), backdrop)
}

func (c *client) Provision(containerID string) error {
	_, err := c.configClient.Provision(context.Background(), &types.ContainerId{Id: containerID})
	return err
}
