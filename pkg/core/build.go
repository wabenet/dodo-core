package core

import (
	"os"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	log "github.com/hashicorp/go-hclog"
	"github.com/moby/term"
)

func BuildImage(config *api.BuildInfo) (string, error) {
	b, err := GetBuilder(config.Builder)
	if err != nil {
		return "", err
	}

	var height, width uint32

	fd := os.Stdin.Fd()
	if term.IsTerminal(fd) {
		state, err := term.SetRawTerminal(fd)
		if err != nil {
			return "", err
		}

		defer func() {
			if err := term.RestoreTerminal(fd, state); err != nil {
				log.L().Error("could not restore terminal", "error", err)
			}
		}()

		ws, err := term.GetWinsize(fd)
		if err != nil {
			return "", err
		}

		height = uint32(ws.Height)
		width = uint32(ws.Width)
	}

	return b.CreateImage(config, &plugin.StreamConfig{
		Stdin:          os.Stdin,
		Stdout:         os.Stdout,
		Stderr:         os.Stderr,
		TerminalHeight: height,
		TerminalWidth:  width,
	})
}
