package proxycmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/command"
	"github.com/dodo-cli/dodo-core/pkg/plugin/runtime"
	"github.com/spf13/cobra"
)

func New(m plugin.Manager, defaultCmd string) *Command {
	cmd := &cobra.Command{
		Use:                name,
		SilenceUsage:       true,
		DisableFlagParsing: true,
		Args:               cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				if self, err := os.Executable(); err == nil {
					return runProxy(self, []string{defaultCmd})
				}

				return fmt.Errorf("could not run plugin '%s': %w", defaultCmd, plugin.ErrPluginNotFound)
			}

			if path, err := exec.LookPath(fmt.Sprintf("dodo-%s", args[0])); err == nil {
				return runProxy(path, args[1:])
			}

			path := plugin.PathByName(args[0])
			if stat, err := os.Stat(path); err == nil && stat.Mode().Perm()&0111 != 0 {
				return runProxy(path, args[1:])
			}

			if self, err := os.Executable(); err == nil {
				return runProxy(self, append([]string{defaultCmd}, args...))
			}

			return fmt.Errorf("could not run plugin '%s': %w", defaultCmd, plugin.ErrPluginNotFound)
		},
	}

	for _, p := range m.GetPlugins(command.Type.String()) {
		cmd.AddCommand(p.(command.Command).GetCobraCommand())
	}

	return &Command{cmd: cmd}
}

func runProxy(executable string, args []string) error {
	cmd := exec.Command(executable, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return &runtime.Result{
				ExitCode: int64(exit.ExitCode()),
				Message:  string(exit.Stderr),
			}
		}

		return err
	}

	return nil
}
