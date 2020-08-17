package run

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dodo/dodo-core/pkg/plugin"
	"github.com/dodo/dodo-core/pkg/plugin/configuration"
	"github.com/dodo/dodo-core/pkg/plugin/runtime"
	"github.com/dodo/dodo-core/pkg/types"
	log "github.com/hashicorp/go-hclog"
	"github.com/moby/term"
	"golang.org/x/net/context"
)

type Container struct {
	daemon  bool
	config  *types.Backdrop
	context context.Context
}

func GetRuntime() (runtime.ContainerRuntime, error) {
	for _, p := range plugin.GetPlugins(runtime.Type.String()) {
		if rt, ok := p.(runtime.ContainerRuntime); ok {
			return rt, nil
		}
	}

	return nil, fmt.Errorf("could not find container runtime: %w", plugin.ErrNoValidPluginFound)
}

func NewContainer(overrides *types.Backdrop, daemon bool) (*Container, error) {
	c := &Container{
		daemon: daemon,
		config: &types.Backdrop{
			Name:       overrides.Name,
			Entrypoint: &types.Entrypoint{},
		},
		context: context.Background(),
	}

	for _, p := range plugin.GetPlugins(configuration.Type.String()) {
		conf, err := p.(configuration.Configuration).UpdateConfiguration(c.config)
		if err != nil {
			log.Default().Warn("could not get config", "error", err)
			continue
		}

		c.config.Merge(conf)
	}

	c.config.Merge(overrides)
	log.Default().Debug("assembled configuration", "backdrop", c.config)

	if c.daemon {
		c.config.ContainerName = c.config.Name
	} else if len(c.config.ContainerName) == 0 {
		id := make([]byte, 8)
		if _, err := rand.Read(id); err != nil {
			panic(err)
		}

		c.config.ContainerName = fmt.Sprintf("%s-%s", c.config.Name, hex.EncodeToString(id))
	}

	return c, nil
}

func (c *Container) Run() error {
	rt, err := GetRuntime()
	if err != nil {
		return err
	}

	imageID, err := rt.ResolveImage(c.config.ImageId)
	if err != nil {
		return err
	}

	c.config.ImageId = imageID

	tty := term.IsTerminal(os.Stdin.Fd()) && term.IsTerminal(os.Stdout.Fd())

	containerID, err := rt.CreateContainer(c.config, !c.daemon && tty, !c.daemon)
	if err != nil {
		return err
	}

	for _, p := range plugin.GetPlugins(configuration.Type.String()) {
		err := p.(configuration.Configuration).Provision(containerID)
		if err != nil {
			log.Default().Warn("could not provision", "error", err)
		}
	}

	if c.daemon {
		return rt.StartContainer(containerID)
	}

	fd := os.Stdin.Fd()
	if term.IsTerminal(fd) {
		state, err := term.SetRawTerminal(fd)
		if err != nil {
			return err
		}

		defer func() {
			if err := term.RestoreTerminal(fd, state); err != nil {
				log.L().Error("could not restore terminal", "error", err)
			}
		}()

		resize(fd, rt, containerID)

		resizeChannel := make(chan os.Signal, 1)
		signal.Notify(resizeChannel, syscall.SIGWINCH)

		go func() {
			for range resizeChannel {
				resize(fd, rt, containerID)
			}
		}()
	}

	return rt.StreamContainer(containerID, os.Stdin, os.Stdout)
}

func (c *Container) Stop() error {
	rt, err := GetRuntime()
	if err != nil {
		return err
	}

	return rt.RemoveContainer(c.config.ContainerName)
}

func resize(fd uintptr, rt runtime.ContainerRuntime, containerID string) {
	ws, err := term.GetWinsize(fd)
	if err != nil {
		return
	}

	height := uint32(ws.Height)
	width := uint32(ws.Width)

	if height == 0 && width == 0 {
		return
	}

	if err := rt.ResizeContainer(containerID, height, width); err != nil {
		log.L().Warn("could not resize terminar", "error", err)
	}
}
