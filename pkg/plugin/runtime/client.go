package runtime

import (
	"io"
	"net"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/hashicorp/go-hclog"
	"golang.org/x/net/context"
)

var _ ContainerRuntime = &client{}

type client struct {
	runtimeClient api.RuntimePluginClient
}

func (c *client) Type() plugin.Type {
	return Type
}

func (c *client) Init() error {
	_, err := c.runtimeClient.Init(context.Background(), &empty.Empty{})
	return err
}

func (c *client) PluginInfo() (*api.PluginInfo, error) {
	return c.runtimeClient.GetPluginInfo(context.Background(), &empty.Empty{})
}

func (c *client) ResolveImage(spec string) (string, error) {
	img, err := c.runtimeClient.GetImage(context.Background(), &api.GetImageRequest{ImageSpec: spec})
	if err != nil {
		return "", err
	}

	return img.ImageId, nil
}

func (c *client) CreateContainer(config *api.Backdrop, tty bool, stdio bool) (string, error) {
	resp, err := c.runtimeClient.CreateContainer(context.Background(), &api.CreateContainerRequest{
		Config: config,
		Tty:    tty,
		Stdio:  stdio,
	})
	if err != nil {
		return "", err
	}

	return resp.ContainerId, nil
}

func (c *client) StartContainer(id string) error {
	_, err := c.runtimeClient.StartContainer(context.Background(), &api.StartContainerRequest{ContainerId: id})
	return err
}

func (c *client) DeleteContainer(id string) error {
	_, err := c.runtimeClient.DeleteContainer(context.Background(), &api.DeleteContainerRequest{ContainerId: id})
	return err
}

func (c *client) ResizeContainer(id string, height uint32, width uint32) error {
	_, err := c.runtimeClient.ResizeContainer(
		context.Background(),
		&api.ResizeContainerRequest{ContainerId: id, Height: height, Width: width},
	)

	return err
}

func (c *client) StreamContainer(id string, r io.Reader, w io.Writer, height uint32, width uint32) error {
	ctx := context.Background()

	connInfo, err := c.runtimeClient.GetStreamingConnection(ctx, &api.GetStreamingConnectionRequest{ContainerId: id})
	if err != nil {
		return err
	}

	addr, err := net.ResolveTCPAddr("tcp", connInfo.Url)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	defer func() {
		if err := conn.CloseWrite(); err != nil {
			log.L().Warn("could not close streaming connection", "error", err)
		}
	}()

	go func() {
		if _, err := io.Copy(w, conn); err != nil {
			log.L().Warn("error reading from streaming connection", "error", err)
		}

		if err := conn.CloseWrite(); err != nil {
			log.L().Warn("could not close streaming connection", "error", err)
		}
	}()

	go func() {
		if _, err := io.Copy(conn, r); err != nil {
			log.L().Warn("error writing to streaming connection", "error", err)
		}
	}()

	result, err := c.runtimeClient.StreamContainer(ctx, &api.StreamContainerRequest{ContainerId: id, Height: height, Width: width})
	if err != nil {
		return err
	}

	return &Result{
		ExitCode: result.ExitCode,
		Message:  result.Message,
	}
}
