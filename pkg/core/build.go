package core

import (
	"fmt"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/builder"
	"github.com/dodo-cli/dodo-core/pkg/ui"
)

func BuildImage(m plugin.Manager, config *api.BuildInfo) (string, error) {
	b, err := builder.GetByName(m, config.Builder)
	if err != nil {
		return "", fmt.Errorf("could not find build plugin for %s: %w", config.Builder, err)
	}

	imageID := ""

	err = ui.NewTerminal().RunInRaw(
		func(t *ui.Terminal) error {
			if id, err := b.CreateImage(config, &plugin.StreamConfig{
				Stdin:          t.Stdin,
				Stdout:         t.Stdout,
				Stderr:         t.Stderr,
				TerminalHeight: t.Height,
				TerminalWidth:  t.Width,
			}); err != nil {
				return fmt.Errorf("error in container I/O stream: %w", err)
			} else {
				imageID = id
			}

			return nil
		},
	)

	return imageID, err
}
