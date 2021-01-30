package run

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/configuration"
	"github.com/dodo-cli/dodo-core/pkg/plugin/runtime"
	"github.com/dodo-cli/dodo-core/pkg/types"
	log "github.com/hashicorp/go-hclog"
	"github.com/moby/term"
)

func RunContainer(overrides *api.Backdrop) error {
	config := GetConfig(overrides)

	if len(config.ContainerName) == 0 {
		id := make([]byte, 8)
		if _, err := rand.Read(id); err != nil {
			panic(err)
		}

		config.ContainerName = fmt.Sprintf("%s-%s", config.Name, hex.EncodeToString(id))
	}

	rt, err := GetRuntime(config.Runtime)
	if err != nil {
		return err
	}

	imageID, err := rt.ResolveImage(config.ImageId)
	if err != nil {
		return err
	}

	config.ImageId = imageID

	tty := term.IsTerminal(os.Stdin.Fd()) && term.IsTerminal(os.Stdout.Fd())

	containerID, err := rt.CreateContainer(config, tty, true)
	if err != nil {
		return err
	}

	var height, width uint32

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

		ws, err := term.GetWinsize(fd)
		if err != nil {
			return err
		}

		height = uint32(ws.Height)
		width = uint32(ws.Width)

		resizeChannel := make(chan os.Signal, 1)
		signal.Notify(resizeChannel, syscall.SIGWINCH)

		go func() {
			for range resizeChannel {
				resize(fd, rt, containerID)
			}
		}()
	}

	return rt.StreamContainer(containerID, os.Stdin, os.Stdout, height, width)
}

func GetRuntime(name string) (runtime.ContainerRuntime, error) {
	for n, p := range plugin.GetPlugins(runtime.Type.String()) {
		if name != "" && name != n {
			continue
		}

		if rt, ok := p.(runtime.ContainerRuntime); ok {
			return rt, nil
		}
	}

	return nil, fmt.Errorf("could not find container runtime: %w", plugin.ErrNoValidPluginFound)
}

func GetConfig(overrides *api.Backdrop) *api.Backdrop {
	config := &api.Backdrop{Name: overrides.Name, Entrypoint: &api.Entrypoint{}}

	for _, p := range plugin.GetPlugins(configuration.Type.String()) {
		info, err := p.PluginInfo()
		if err != nil {
			log.L().Warn("could not read plugin info")
			continue
		}

		log.L().Debug("Fetching configuration from plugin", "name", info.Name)

		conf, err := p.(configuration.Configuration).GetBackdrop(config.Name)
		if err != nil {
			log.L().Warn("could not get config", "error", err)
			continue
		}

		types.Merge(config, conf)
	}

	types.Merge(config, overrides)
	log.L().Debug("assembled configuration", "backdrop", config)

	return config
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
		log.L().Warn("could not resize terminal", "error", err)
	}
}
