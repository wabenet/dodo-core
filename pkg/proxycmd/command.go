package proxycmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dodo/dodo-core/pkg/plugin"
	"github.com/dodo/dodo-core/pkg/plugin/command"
	"github.com/dodo/dodo-core/pkg/types"
	"github.com/spf13/cobra"
)

func Execute(defaultCmd string) int {
	cmd := &cobra.Command{
		Use:                "dodo",
		SilenceUsage:       true,
		DisableFlagParsing: true,
		Args:               cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				if self, err := os.Executable(); err == nil {
					return runProxy(self, []string{defaultCmd})
				} else {
					return fmt.Errorf("could not run plugin '%s': %w", defaultCmd, plugin.ErrPluginNotFound)
				}
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

	for _, p := range plugin.GetPlugins(command.Type.String()) {
		cmd.AddCommand(p.(command.Command).GetCobraCommand())
	}

	if err := cmd.Execute(); err != nil {
		if err, ok := err.(*types.Result); ok {
			return int(err.ExitCode)
		}

		return 1
	}

	return 0
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
