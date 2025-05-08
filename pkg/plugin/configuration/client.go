package configuration

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	api "github.com/wabenet/dodo-core/api/configuration/v1alpha2"
	pluginapi "github.com/wabenet/dodo-core/api/plugin/v1alpha1"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"google.golang.org/grpc"
)

var _ Configuration = &Client{}

type Client struct {
	configClient api.PluginClient
}

func NewGRPCClient(conn grpc.ClientConnInterface) Configuration {
	return &Client{configClient: api.NewPluginClient(conn)}
}

func (c *Client) Type() plugin.Type { //nolint:ireturn
	return Type
}

func (c *Client) PluginInfo() *pluginapi.PluginInfo {
	info, err := c.configClient.GetPluginInfo(context.Background(), &empty.Empty{})
	if err != nil {
		pi := &pluginapi.PluginInfo{}

		pi.SetName(plugin.MkName(Type, plugin.FailedPlugin))
		pi.SetFields(map[string]string{"error": err.Error()})

		return pi
	}

	return info
}

func (c *Client) Init() (plugin.Config, error) {
	resp, err := c.configClient.InitPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("could not initialize plugin: %w", err)
	}

	return resp.GetConfig(), nil
}

func (c *Client) Cleanup() {
	_, err := c.configClient.ResetPlugin(context.Background(), &empty.Empty{})
	if err != nil {
		log.L().Error("plugin reset error", "error", err)
	}
}

func (c *Client) ListBackdrops() ([]Backdrop, error) {
	response, err := c.configClient.ListBackdrops(context.Background(), &empty.Empty{})
	if err != nil {
		return []Backdrop{}, fmt.Errorf("could not list backdrops: %w", err)
	}

	result := []Backdrop{}

	for _, b := range response.GetBackdrops() {
		result = append(result, BackdropFromProto(b))
	}

	return result, nil
}

func (c *Client) GetBackdrop(alias string) (Backdrop, error) {
	req := &api.GetBackdropRequest{}

	req.SetAlias(alias)

	response, err := c.configClient.GetBackdrop(context.Background(), req)
	if err != nil {
		return EmptyBackdrop(), fmt.Errorf("could not get backdrop: %w", err)
	}

	return BackdropFromProto(response.GetBackdrop()), nil
}
