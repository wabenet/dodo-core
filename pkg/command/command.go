package command

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/hashicorp/go-hclog"
	"github.com/oclaussen/dodo/pkg/appconfig"
	"github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/dodo/pkg/plugin/command"
	"github.com/oclaussen/dodo/pkg/plugin/configuration"
	"github.com/oclaussen/dodo/pkg/plugin/runtime"
	"github.com/oclaussen/dodo/pkg/run"
	"github.com/oclaussen/dodo/pkg/types"
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

func Execute() int {
	log.SetDefault(log.New(appconfig.GetLoggerOptions()))

	plugin.RegisterBuiltin(&run.Command{})

	plugin.LoadPlugins(
		command.Type,
		configuration.Type,
		runtime.Type,
	)
	defer plugin.UnloadPlugins()

	if err := NewCommand().Execute(); err != nil {
		if err, ok := err.(*types.Result); ok {
			return int(err.ExitCode)
		}

		return 1
	}

	return 0
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "dodo",
		Short:              "Run commands in a Docker context",
		Long:               description,
		SilenceUsage:       true,
		DisableFlagParsing: true,
		Args:               cobra.MinimumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if path, err := exec.LookPath(fmt.Sprintf("dodo-%s", args[0])); err == nil {
				return runProxy(path, args[1:])
			}

			path := plugin.PathByName(args[0])
			if stat, err := os.Stat(path); err == nil && stat.Mode().Perm()&0111 != 0 {
				return runProxy(path, args[1:])
			}

			if self, err := os.Executable(); err == nil {
				return runProxy(self, append([]string{"run"}, args...))
			}

			return fmt.Errorf("could not run: %w", plugin.ErrPluginNotFound)
		},
	}

	for _, p := range plugin.GetPlugins(command.Type.String()) {
		cmd.AddCommand(p.(command.Command).GetCobraCommand())
	}

	return cmd
}

func runProxy(executable string, args []string) error {
	cmd := exec.Command(executable, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return &types.Result{
				ExitCode: int64(exit.ExitCode()),
				Message:  string(exit.Stderr),
			}
		}

		return err
	}

	return nil
}
