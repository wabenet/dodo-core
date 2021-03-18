package builder

import (
	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

var _ ImageBuilder = &client{}

type client struct {
	builderClient api.BuilderPluginClient
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) PluginInfo() (*api.PluginInfo, error) {
	return c.builderClient.GetPluginInfo(context.Background(), &empty.Empty{})
}

func (c *client) CreateImage(config *api.BuildInfo) (string, error) {
	resp, err := c.builderClient.CreateImage(context.Background(), &api.CreateImageRequest{Config: config})
	if err != nil {
		return "", err
	}

	return resp.ImageId, nil
}
