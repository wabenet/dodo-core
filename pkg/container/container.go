package container

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/pkg/stringid"
	"github.com/moby/term"
	"github.com/oclaussen/dodo/pkg/configuration"
	"github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/dodo/pkg/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Container struct {
	daemon  bool
	config  *types.Backdrop
	context context.Context
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

	for _, p := range plugin.GetPlugins(configuration.PluginType) {
		conf, err := p.(configuration.Configuration).UpdateConfiguration(c.config)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("could not get config")
			continue
		}
		c.config = conf
	}
	c.config.Merge(overrides)

	log.WithFields(log.Fields{"backdrop": c.config}).Debug("assembled configuration")

	if c.daemon {
		c.config.ContainerName = c.config.Name
	} else if len(c.config.ContainerName) == 0 {
		c.config.ContainerName = fmt.Sprintf("%s-%s", c.config.Name, stringid.GenerateRandomID()[:8])
	}

	return c, nil
}

func (c *Container) Run() error {
	rt, err := GetRuntime()
	if err != nil {
		return err
	}

	imageId, err := rt.ResolveImage(c.config.ImageId)
	if err != nil {
		return err
	}
	c.config.ImageId = imageId

	containerID, err := rt.CreateContainer(c.config)
	if err != nil {
		return err
	}

	for _, p := range plugin.GetPlugins(configuration.PluginType) {
		err := p.(configuration.Configuration).Provision(containerID)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("could not provision")
		}
	}

	if c.daemon {
		return rt.StartContainer(containerID)
	}

	inFd, inTerm := term.GetFdInfo(os.Stdin)
	outFd, outTerm := term.GetFdInfo(os.Stdout)

	if inTerm && outTerm {
		inState, err := term.SetRawTerminal(inFd)
		if err != nil {
			return err
		}
		defer term.RestoreTerminal(inFd, inState)

		outState, err := term.SetRawTerminal(outFd)
		if err != nil {
			return err
		}
		defer term.RestoreTerminal(outFd, outState)

		resize(rt, containerID)
		resizeChannel := make(chan os.Signal, 1)
		signal.Notify(resizeChannel, syscall.SIGWINCH)
		go func() {
			for range resizeChannel {
				resize(rt, containerID)
			}
		}()
	}

        return rt.StreamContainer(containerID, os.Stdout, os.Stdin)
}

func (c *Container) Stop() error {
	rt, err := GetRuntime()
	if err != nil {
		return err
	}
	return rt.RemoveContainer(c.config.ContainerName)
}

func resize(rt ContainerRuntime, containerID string) {
	outFd, _ := term.GetFdInfo(os.Stdout)

	ws, err := term.GetWinsize(outFd)
	if err != nil {
		return
	}

	height := uint32(ws.Height)
	width := uint32(ws.Width)
	if height == 0 && width == 0 {
		return
	}

	rt.ResizeContainer(containerID, height, width)
}
