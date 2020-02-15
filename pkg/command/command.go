package command

import (
	"github.com/oclaussen/dodo/pkg/config"
	"github.com/oclaussen/dodo/pkg/container"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// TODO: clean up command structure and options

const description = `Run commands in a Docker context.

Dodo operates on a set of backdrops, that must be configured in configuration
files (in the current directory or one of the config directories). Backdrops
are similar to docker-composes services, but they define one-shot commands
instead of long-running services. More specifically, each backdrop defines a 
docker container in which a script should be executed. Dodo simply passes all 
CMD arguments to the first backdrop with NAME that is found. Additional FLAGS
can be used to overwrite the backdrop configuration.
`

func NewCommand() *cobra.Command {
	var opts options
	cmd := &cobra.Command{
		Use:                   "dodo [flags] [name] [cmd...]",
		Short:                 "Run commands in a Docker context",
		Long:                  description,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configureLogging()

			backdrop, err := opts.createConfig(args[0], args[1:])
			if err != nil {
				return err
			}

			c, err := container.NewContainer(backdrop, config.LoadAuthConfig(), false)
			if err != nil {
				return err
			}

			return c.Run()
		},
	}
	opts.createFlags(cmd)

	cmd.AddCommand(NewRunCommand())
	return cmd
}

func NewRunCommand() *cobra.Command {
	var opts options
	cmd := &cobra.Command{
		Use:                   "run [flags] [name] [cmd...]",
		Short:                 "Same as running 'dodo [name]', can be used when a backdrop name collides with a top-level command",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			configureLogging()

			backdrop, err := opts.createConfig(args[0], args[1:])
			if err != nil {
				return err
			}

			c, err := container.NewContainer(backdrop, config.LoadAuthConfig(), false)
			if err != nil {
				return err
			}

			return c.Run()
		},
	}

	opts.createFlags(cmd)
	return cmd
}

func configureLogging() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})
}
