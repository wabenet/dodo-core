package core

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/runtime"
	log "github.com/hashicorp/go-hclog"
	"github.com/moby/term"
	"golang.org/x/sync/errgroup"
)

func RunBackdrop(m plugin.Manager, b *api.Backdrop) (int, error) {
	rt, err := runtime.GetByName(m, b.Runtime)
	if err != nil {
		return ExitCodeInternalError, err
	}

	imageID, err := rt.ResolveImage(b.ImageId)
	if err != nil {
		return ExitCodeInternalError, err
	}

	b.ImageId = imageID

	tty := term.IsTerminal(os.Stdin.Fd()) && term.IsTerminal(os.Stdout.Fd())

	containerID, err := rt.CreateContainer(b, tty, true)
	if err != nil {
		return ExitCodeInternalError, err
	}

	var height, width uint32

	eg, _ := errgroup.WithContext(context.Background())
	resizeChannel := make(chan os.Signal, 1)

	if fd := os.Stdin.Fd(); term.IsTerminal(fd) {
		state, err := term.SetRawTerminal(fd)
		if err != nil {
			return ExitCodeInternalError, fmt.Errorf("could not set raw terminal: %w", err)
		}

		defer func() {
			if err := term.RestoreTerminal(fd, state); err != nil {
				log.L().Error("could not restore terminal", "error", err)
			}
		}()

		ws, err := term.GetWinsize(fd)
		if err != nil {
			return ExitCodeInternalError, fmt.Errorf("could not get terminal size: %w", err)
		}

		height = uint32(ws.Height)
		width = uint32(ws.Width)

		signal.Notify(resizeChannel, syscall.SIGWINCH)

		eg.Go(func() error {
			for range resizeChannel {
				resize(fd, rt, containerID)
			}

			return nil
		})
	}

	exitCode := 0

	eg.Go(func() error {
		defer close(resizeChannel)

		r, err := rt.StreamContainer(containerID, &plugin.StreamConfig{
			Stdin:          os.Stdin,
			Stdout:         os.Stdout,
			Stderr:         os.Stderr,
			TerminalHeight: height,
			TerminalWidth:  width,
		})

		exitCode = r.ExitCode

		return err
	})

	return exitCode, eg.Wait()
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
