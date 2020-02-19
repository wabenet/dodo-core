package command

import (
	"github.com/oclaussen/dodo/pkg/container"
	"github.com/spf13/cobra"
)

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
		RunE: func(_ *cobra.Command, args []string) error {
			return run(&opts, args[0], args[1:])
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
		RunE: func(_ *cobra.Command, args []string) error {
			return run(&opts, args[0], args[1:])
		},
	}

	opts.createFlags(cmd)
	return cmd
}

func run(opts *options, name string, command []string) error {
	backdrop, err := opts.createConfig(name, command)
	if err != nil {
		return err
	}

	c, err := container.NewContainer(backdrop, false)
	if err != nil {
		return err
	}

	return c.Run()
}
