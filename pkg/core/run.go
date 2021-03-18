package core

import (
	"os"
	"os/signal"
	"syscall"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/plugin/runtime"
	log "github.com/hashicorp/go-hclog"
	"github.com/moby/term"
)

func RunBackdrop(config *api.Backdrop) error {
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
