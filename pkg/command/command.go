package command

import (
	"github.com/docker/docker/client"
	"github.com/oclaussen/dodo/pkg/config"
	"github.com/oclaussen/dodo/pkg/container"
	"github.com/oclaussen/dodo/pkg/image"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// TODO: add some --no-rm option?
// TODO: missing environment, user, volumes, volumes_from
// TODO: go through options of docker, docker-compose and sudo
type options struct {
	Filename    string
	Debug       bool
	Interactive bool
	Remove      bool
	NoCache     bool
	Pull        bool
	Build       bool
	Workdir     string
}

// TODO: no error message when bind mount fails

// TODO: tests

// NewCommand creates a new command instance
func NewCommand() *cobra.Command {
	var opts options
	cmd := &cobra.Command{
		Use:              "dodo [OPTIONS] NAME [CMD...]",
		Short:            "Run commands in a Docker context",
		Version:          "0.0.1", // TODO: fix help/version/errors
		TraverseChildren: true,
		Args:             cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(&opts, args[0], args[1:])
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Filename, "file", "f", "", "Specify a dodo configuration file")
	flags.BoolVarP(&opts.Debug, "debug", "", false, "Show additional debug output")
	flags.BoolVarP(&opts.Interactive, "interactive", "i", false, "Run an interactive session")
	flags.BoolVarP(&opts.Remove, "rm", "", false, "Automatically remove the container when it exits")
	flags.BoolVarP(&opts.NoCache, "no-cache", "", false, "Do not use cache when building the image")
	flags.BoolVarP(&opts.Pull, "pull", "", false, "Always attempt to pull a newer version of the image")
	flags.BoolVarP(&opts.Build, "build", "", false, "Always build an image, even if already exists")
	flags.StringVarP(&opts.Workdir, "workdir", "w", "", "Working directory inside the container")
	flags.SetInterspersed(false)

	return cmd
}

func runCommand(opts *options, name string, command []string) error {
	config, err := config.LoadConfiguration(name, opts.Filename)
	if err != nil {
		return err
	}

	if len(command) > 0 {
		config.Command = command
	}
	if opts.Debug {
		log.SetLevel(log.DebugLevel)
	}
	if opts.Remove {
		remove := true
		config.Remove = &remove
	}
	if opts.Workdir != "" {
		// TODO: this does not seem to work?
		config.WorkingDir = opts.Workdir
	}
	if opts.Interactive {
		config.Interactive = true
	}
	if opts.Pull {
		config.Pull = true
	}
	if config.Build != nil && opts.NoCache {
		config.Build.NoCache = true
	}
	if config.Build != nil && opts.Build {
		config.Build.ForceRebuild = true
	}

	// TODO: read docker configuration
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	ctx := context.Background()

	imageID, err := image.Get(ctx, image.Options{
		Client:    dockerClient,
		Name:      config.Image,
		Build:     config.Build,
		ForcePull: config.Pull,
	})
	if err != nil {
		return err
	}

	// TODO: generate a temp file in the container for the entrypoint
	// TODO feels inefficient to stupid all of config
	return container.Run(ctx, container.Options{
		Client:      dockerClient,
		Image:       imageID,
		Name:        config.ContainerName,
		Interactive: config.Interactive,
		Interpreter: config.Interpreter,
		Entrypoint:  "/tmp/dodo-entrypoint",
		Script:      config.Script,
		Command:     config.Command,
		Environment: config.Environment,
		Volumes:     config.Volumes,
		VolumesFrom: config.VolumesFrom,
		User:        config.User,
		WorkingDir:  config.WorkingDir,
	})
}
