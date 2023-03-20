package configuration

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	configuration "github.com/wabenet/dodo-core/api/configuration/v1alpha1"
	core "github.com/wabenet/dodo-core/api/core/v1alpha5"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"google.golang.org/grpc"
)

var _ Configuration = &client{}

type client struct {
	configClient configuration.PluginClient
}

func NewGRPCClient(conn grpc.ClientConnInterface) Configuration {
	return &client{configClient: configuration.NewPluginClient(conn)}
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) PluginInfo() *core.PluginInfo {
	info, err := c.configClient.GetPluginInfo(context.Background(), &empty.Empty{})
	if err != nil {
		return &core.PluginInfo{
			Name:   &core.PluginName{Type: Type.String(), Name: plugin.FailedPlugin},
			Fields: map[string]string{"error": err.Error()},
		}
	}

	return info
}

func (c *client) Init() (plugin.Config, error) {
	resp, err := c.configClient.InitPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return resp.Config, nil
}

func (c *client) Cleanup() {
	_, err := c.configClient.ResetPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		log.L().Error("plugin reset error", "error", err)
	}
}

func (c *client) ListBackdrops() ([]*core.Backdrop, error) {
	response, err := c.configClient.ListBackdrops(context.Background(), &empty.Empty{})
	if err != nil {
		return []*core.Backdrop{}, fmt.Errorf("could not list backdrops: %w", err)
	}

	return response.Backdrops, nil
}

func (c *client) GetBackdrop(alias string) (*core.Backdrop, error) {
	response, err := c.configClient.GetBackdrop(context.Background(), &configuration.GetBackdropRequest{Alias: alias})
	if err != nil {
		return nil, fmt.Errorf("could not get backdrop: %w", err)
	}

	return response, nil
}
