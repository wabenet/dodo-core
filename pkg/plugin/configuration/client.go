package configuration

import (
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/types"
	"golang.org/x/net/context"
)

var _ Configuration = &client{}

type client struct {
	configClient types.ConfigurationClient
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) Init() (*types.PluginInfo, error) {
	return c.configClient.Init(context.Background(), &types.Empty{})
}

func (c *client) UpdateConfiguration(backdrop *types.Backdrop) (*types.Backdrop, error) {
	return c.configClient.UpdateConfiguration(context.Background(), backdrop)
}

func (c *client) Provision(containerID string) error {
	_, err := c.configClient.Provision(context.Background(), &types.ContainerId{Id: containerID})
	return err
}
