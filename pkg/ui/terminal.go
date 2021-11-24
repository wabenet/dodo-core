package ui

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	log "github.com/hashicorp/go-hclog"
	"github.com/moby/term"
	"golang.org/x/sync/errgroup"
)

type Terminal struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Height uint32
	Width  uint32

	resizeHook    func(*Terminal)
	resizeChannel chan os.Signal
}

func NewTerminal() *Terminal {
	return &Terminal{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		resizeHook: nil,
	}
}

func (t *Terminal) WithResizeHook(f func(*Terminal)) *Terminal {
	t.resizeHook = f
	t.resizeChannel = make(chan os.Signal, 1)

	return t
}

func (t *Terminal) RunInRaw(wrapped func(*Terminal) error) error {
	eg, _ := errgroup.WithContext(context.Background())

	if fd := os.Stdin.Fd(); term.IsTerminal(fd) {
		state, err := term.SetRawTerminal(fd)
		if err != nil {
			return fmt.Errorf("could not set raw terminal: %w", err)
		}

		defer func() {
			if err := term.RestoreTerminal(fd, state); err != nil {
				log.L().Error("could not restore terminal", "error", err)
			}
		}()

		ws, err := term.GetWinsize(fd)
		if err != nil {
			return fmt.Errorf("could not get terminal size: %w", err)
		}

		t.Height = uint32(ws.Height)
		t.Width = uint32(ws.Width)

		if t.resizeHook != nil {
			signal.Notify(t.resizeChannel, syscall.SIGWINCH)

			eg.Go(func() error {
				for range t.resizeChannel {
					ws, err := term.GetWinsize(fd)
					if err != nil {
						continue
					}

					t.Height = uint32(ws.Height)
					t.Width = uint32(ws.Width)

					if t.Height > 0 && t.Width > 0 {
						t.resizeHook(t)
					}
				}

				return nil
			})
		}
	}

	eg.Go(func() error {
		if t.resizeHook != nil {
			defer close(t.resizeChannel)
		}

		return wrapped(t)
	})

	return eg.Wait()
}

func IsTTY() bool {
	return term.IsTerminal(os.Stdin.Fd()) && term.IsTerminal(os.Stdout.Fd())
}
