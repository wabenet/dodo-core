package run

import (
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/appconfig"
	"github.com/dodo-cli/dodo-core/pkg/core"
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

type options struct {
	interactive bool
	user        string
	workdir     string
	volumes     []string
	environment []string
	publish     []string
	runtime     string
}

func NewCommand() *cobra.Command {
	var opts options

	cmd := &cobra.Command{
		Use:                   name,
		Short:                 "Run commands in Docker context",
		Long:                  description,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		// Because exec results are not really errors
		// TODO: is there a more elegant way to handle this?
		SilenceErrors: true,
		Args:          cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			backdrop, err := opts.createConfig(args[0], args[1:])
			if err != nil {
				return err
			}

			return core.RunByName(backdrop)
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
	flags.StringVarP(
		&opts.runtime, "runtime", "r", "",
		"select runtime plugin")

	return cmd
}

func (opts *options) createConfig(name string, command []string) (*api.Backdrop, error) {
	config := &api.Backdrop{
		Name:    name,
		Runtime: opts.runtime,
		Entrypoint: &api.Entrypoint{
			Interactive: opts.interactive,
			Arguments:   command,
		},
		User:       opts.user,
		WorkingDir: opts.workdir,
	}

	for _, spec := range opts.volumes {
		vol, err := appconfig.ParseVolumeMount(spec)
		if err != nil {
			return nil, err
		}

		config.Volumes = append(config.Volumes, vol)
	}

	for _, spec := range opts.environment {
		env, err := appconfig.ParseEnvironmentVariable(spec)
		if err != nil {
			return nil, err
		}

		config.Environment = append(config.Environment, env)
	}

	for _, spec := range opts.publish {
		port, err := appconfig.ParsePortBinding(spec)
		if err != nil {
			return nil, err
		}

		config.Ports = append(config.Ports, port)
	}

	return config, nil
}
