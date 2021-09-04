package build

import (
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/core"
	"github.com/spf13/cobra"
)

type options struct {
	noCache      bool
	forceRebuild bool
	forcePull    bool
}

func NewCommand() *cobra.Command {
	var opts options

	cmd := &cobra.Command{
		Use:                   name,
		Short:                 "Build all required images for backdrop without running it",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := &api.BuildInfo{
				ImageName:    args[0],
				NoCache:      opts.noCache,
				ForceRebuild: opts.forceRebuild,
				ForcePull:    opts.forcePull,
			}

			_, err := core.BuildByName(config)

			return err
		},
	}

	flags := cmd.Flags()
	flags.SetInterspersed(false)

	flags.BoolVar(
		&opts.noCache, "no-cache", false,
		"do not use cache when building the image")
	flags.BoolVarP(
		&opts.forceRebuild, "force", "f", false,
		"always rebuild all dependencies, even when they already exist")
	flags.BoolVar(
		&opts.forcePull, "pull", false,
		"always attempt to pull base images")

	return cmd
}
