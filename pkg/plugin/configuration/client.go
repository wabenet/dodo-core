package configuration

import (
	"fmt"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

var _ Configuration = &client{}

type client struct {
	configClient api.ConfigurationPluginClient
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) PluginInfo() (*api.PluginInfo, error) {
	return c.configClient.GetPluginInfo(context.Background(), &empty.Empty{})
}

func (c *client) ListBackdrops() ([]*api.Backdrop, error) {
	response, err := c.configClient.ListBackdrops(context.Background(), &empty.Empty{})
	if err != nil {
		return []*api.Backdrop{}, fmt.Errorf("could not list backdrops: %w", err)
	}

	return response.Backdrops, nil
}

func (c *client) GetBackdrop(alias string) (*api.Backdrop, error) {
	return c.configClient.GetBackdrop(context.Background(), &api.GetBackdropRequest{Alias: alias})
}
