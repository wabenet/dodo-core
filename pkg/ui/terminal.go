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

	fd uintptr

	signalChannel chan os.Signal
	signalHook    func(os.Signal, *Terminal)
}

func NewTerminal() *Terminal {
	return &Terminal{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,

		fd: os.Stdin.Fd(),

		signalChannel: make(chan os.Signal, 128),
		signalHook:    nil,
	}
}

func (t *Terminal) OnSignal(f func(os.Signal, *Terminal)) *Terminal {
	t.signalHook = f

	return t
}

func (t *Terminal) refreshSize() {
	if !term.IsTerminal(t.fd) {
		return
	}

	ws, err := term.GetWinsize(t.fd)
	if err != nil {
		log.L().Warn("could not get terminal size", "error", err)

		return
	}

	t.Height = uint32(ws.Height)
	t.Width = uint32(ws.Width)
}

func (t *Terminal) stopSignals() {
	signal.Stop(t.signalChannel)
	close(t.signalChannel)
}

func (t *Terminal) handleSignal(s os.Signal) {
	// Special cases handled here
	switch s {
	case syscall.SIGCHLD, syscall.SIGPIPE:
		return

	case syscall.SIGWINCH:
		t.refreshSize()
	}

	t.signalHook(s, t)
}

func (t *Terminal) RunInRaw(wrapped func(*Terminal) error) error {
	eg, _ := errgroup.WithContext(context.Background())

	if term.IsTerminal(t.fd) {
		state, err := term.SetRawTerminal(t.fd)
		if err != nil {
			return fmt.Errorf("could not set raw terminal: %w", err)
		}

		defer func() {
			if err := term.RestoreTerminal(t.fd, state); err != nil {
				log.L().Error("could not restore terminal", "error", err)
			}
		}()

		t.refreshSize()
	}

	signal.Notify(t.signalChannel)

	eg.Go(func() error {
		for s := range t.signalChannel {
			t.handleSignal(s)
		}

		return nil
	})

	eg.Go(func() error {
		defer t.stopSignals()

		return wrapped(t)
	})

	return eg.Wait()
}

func IsTTY() bool {
	return term.IsTerminal(os.Stdin.Fd()) && term.IsTerminal(os.Stdout.Fd())
}
