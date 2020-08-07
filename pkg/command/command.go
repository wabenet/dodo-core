package command

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	log "github.com/hashicorp/go-hclog"
	"github.com/oclaussen/dodo/pkg/appconfig"
	"github.com/oclaussen/dodo/pkg/plugin"
	configplugin "github.com/oclaussen/dodo/pkg/plugin/configuration"
	runtimeplugin "github.com/oclaussen/dodo/pkg/plugin/runtime"
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

var builtinPlugins = map[string]*cobra.Command{}

func Execute() int {
	log.SetDefault(log.New(appconfig.GetLoggerOptions()))

	plugin.LoadPlugins(
		configplugin.Type,
		runtimeplugin.Type,
	)
	defer plugin.UnloadPlugins()

	Register("run", run.NewCommand())

	if err := NewCommand().Execute(); err != nil {
		if err, ok := err.(*types.Result); ok {
			return int(err.ExitCode)
		}

		return 1
	}

	return 0
}

func Register(name string, cmd *cobra.Command) {
	builtinPlugins[name] = cmd
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
			executable, execArgs, err := findPluginExecutable(args[0])
			if err == nil {
				return runPlugin(executable, execArgs, args[1:])
			}

			executable, execArgs, err = findPluginExecutable("run")
			if err == nil {
				return runPlugin(executable, execArgs, args)
			}
			return err
		},
	}

	for _, subCmd := range builtinPlugins {
		cmd.AddCommand(subCmd)
	}

	return cmd
}

func runPlugin(executable string, execArgs []string, args []string) error {
	cmd := exec.Command(executable, append(execArgs, args...)...)

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

func findPluginExecutable(name string) (string, []string, error) {
	if _, ok := builtinPlugins[name]; ok {
		if self, err := os.Executable(); err == nil {
			return self, []string{name}, nil
		}
	}

	nameInPath := fmt.Sprintf("dodo-%s", name)
	if plugin, err := exec.LookPath(nameInPath); err == nil {
		return plugin, []string{}, nil
	}

	nameInPlugins := filepath.Join(
		appconfig.GetPluginDir(),
		fmt.Sprintf("dodo-%s_%s_%s", name, runtime.GOOS, runtime.GOARCH),
	)
	if stat, err := os.Stat(nameInPlugins); err == nil && stat.Mode().Perm()&0111 != 0 {
		return nameInPlugins, []string{}, nil
	}

	return "", []string{}, fmt.Errorf("%s: %w", name, plugin.ErrPluginNotFound)
}
