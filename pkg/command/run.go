package command

import (
	"github.com/oclaussen/dodo/pkg/configuration"
	"github.com/oclaussen/dodo/pkg/container"
	"github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/dodo/pkg/types"
	"github.com/spf13/cobra"
)

type options struct {
	interactive bool
	user        string
	workdir     string
	volumes     []string
	environment []string
	publish     []string
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
			container.RegisterPlugin()
			configuration.RegisterPlugin()

			plugin.LoadPlugins()
			defer plugin.UnloadPlugins()

			backdrop, err := opts.createConfig(args[0], args[1:])
			if err != nil {
				return err
			}

			c, err := container.NewContainer(backdrop, false)
			if err != nil {
				return err
			}

			return c.Run()
		},
	}

	flags := cmd.Flags()
	flags.SetInterspersed(false)

	flags.BoolVarP(
		&opts.interactive, "interactive", "i", false,
		"run an interactive session")
	flags.StringVarP(
		&opts.user, "user", "u", "",
		"username or UID (format: <name|uid>[:<group|gid>])")
	flags.StringVarP(
		&opts.workdir, "workdir", "w", "",
		"working directory inside the container")
	flags.StringArrayVarP(
		&opts.volumes, "volume", "v", []string{},
		"bind mount a volume")
	flags.StringArrayVarP(
		&opts.environment, "env", "e", []string{},
		"set environment variables")
	flags.StringArrayVarP(
		&opts.publish, "publish", "p", []string{},
		"publish a container's port(s) to the host")

	return cmd
}

func (opts *options) createConfig(name string, command []string) (*types.Backdrop, error) {
	config := &types.Backdrop{
		Name: name,
		Entrypoint: &types.Entrypoint{
			Interactive: opts.interactive,
			Arguments:   command,
		},
		User:       opts.user,
		WorkingDir: opts.workdir,
	}

	for _, spec := range opts.volumes {
		vol := &types.Volume{}
		if err := vol.FromString(spec); err != nil {
			return nil, err
		}

		config.Volumes = append(config.Volumes, vol)
	}

	for _, spec := range opts.environment {
		env := &types.Environment{}
		if err := env.FromString(spec); err != nil {
			return nil, err
		}

		config.Environment = append(config.Environment, env)
	}

	for _, spec := range opts.publish {
		port := &types.Port{}
		if err := port.FromString(spec); err != nil {
			return nil, err
		}

		config.Ports = append(config.Ports, port)
	}

	return config, nil
}
