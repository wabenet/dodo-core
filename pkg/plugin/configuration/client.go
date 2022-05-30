package configuration

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/wabenet/dodo-core/api/v1alpha3"
	"github.com/wabenet/dodo-core/pkg/plugin"
)

var _ Configuration = &client{}

type client struct {
	configClient api.ConfigurationPluginClient
}

func NewGRPCClient(c api.ConfigurationPluginClient) Configuration {
	return &client{configClient: c}
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) PluginInfo() *api.PluginInfo {
	info, err := c.configClient.GetPluginInfo(context.Background(), &empty.Empty{})
	if err != nil {
		return &api.PluginInfo{
			Name:   &api.PluginName{Type: Type.String(), Name: plugin.FailedPlugin},
			Fields: map[string]string{"error": err.Error()},
		}
	}

	return info
}

func (c *client) Init() (plugin.PluginConfig, error) {
	resp, err := c.configClient.InitPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return resp.Config, nil
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
