package runtime

import (
	"io"
	"net"

	log "github.com/hashicorp/go-hclog"
	"github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/dodo/pkg/types"
	"golang.org/x/net/context"
)

type client struct {
	runtimeClient types.ContainerRuntimeClient
}

func (t *client) Type() plugin.Type {
	return Type
}

func (c *client) Init() error {
	_, err := c.runtimeClient.Init(context.Background(), &types.Empty{})
	return err
}

func (c *client) ResolveImage(spec string) (string, error) {
	img, err := c.runtimeClient.ResolveImage(context.Background(), &types.Image{Name: spec})
	if err != nil {
		return "", err
	}

	return img.Id, nil
}

func (c *client) CreateContainer(config *types.Backdrop, tty bool, stdio bool) (string, error) {
	resp, err := c.runtimeClient.CreateContainer(context.Background(), &types.ContainerConfig{
		Config: config,
		Tty:    tty,
		Stdio:  stdio,
	})
	if err != nil {
		return "", err
	}

	return resp.Id, nil
}

func (c *client) StartContainer(id string) error {
	_, err := c.runtimeClient.StartContainer(context.Background(), &types.ContainerId{Id: id})
	return err
}

func (c *client) RemoveContainer(id string) error {
	_, err := c.runtimeClient.RemoveContainer(context.Background(), &types.ContainerId{Id: id})
	return err
}

func (c *client) ResizeContainer(id string, height uint32, width uint32) error {
	_, err := c.runtimeClient.ResizeContainer(
		context.Background(),
		&types.ContainerBox{Id: id, Height: height, Width: width},
	)

	return err
}

func (c *client) StreamContainer(id string, r io.Reader, w io.Writer) error {
	ctx := context.Background()

	connInfo, err := c.runtimeClient.SetupStreamingConnection(ctx, &types.ContainerId{Id: id})
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
	}()

	go func() {
		if _, err := io.Copy(conn, r); err != nil {
			log.L().Warn("error writing to streaming connection", "error", err)
		}
	}()

	result, err := c.runtimeClient.StreamContainer(ctx, &types.ContainerId{Id: id})
	if err != nil {
		return err
	}

	return result
}
