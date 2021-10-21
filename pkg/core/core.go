package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/command/dodo"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/command"
	"github.com/dodo-cli/dodo-core/pkg/plugin/configuration"
)

const (
	ExitCodeInternalError = 1
	DefaultCommand        = "run"
)

func ExecuteDodoMain(m plugin.Manager) int {
	cmd := dodo.New(m, DefaultCommand).GetCobraCommand()

	if err := cmd.Execute(); err != nil {
		return ExitCodeInternalError
	}

	return command.GetExitCode(cmd)
}

func RunByName(m plugin.Manager, overrides *api.Backdrop) (int, error) {
	b := configuration.AssembleBackdropConfig(m, overrides.Name, overrides)

	if len(b.ContainerName) == 0 {
		id := make([]byte, 8)
		if _, err := rand.Read(id); err != nil {
			panic(err)
		}

		b.ContainerName = fmt.Sprintf("%s-%s", b.Name, hex.EncodeToString(id))
	}

	if len(b.ImageId) == 0 {
		if b.BuildInfo == nil {
			return ExitCodeInternalError, fmt.Errorf("neither image nor build configured for backdrop '%s'", overrides.Name)
		}

		for _, dep := range b.BuildInfo.Dependencies {
			if _, err := BuildByName(m, &api.BuildInfo{ImageName: dep}); err != nil {
				return ExitCodeInternalError, err
			}
		}

		imageID, err := BuildImage(m, b.BuildInfo)
		if err != nil {
			return ExitCodeInternalError, err
		}

		b.ImageId = imageID
	}

	return RunBackdrop(m, b)
}

func BuildByName(m plugin.Manager, overrides *api.BuildInfo) (string, error) {
	config, err := configuration.FindBuildConfig(m, overrides.ImageName, overrides)
	if err != nil {
		return "", err
	}

	for _, dep := range config.Dependencies {
		conf := &api.BuildInfo{}
		configuration.MergeBuildInfo(conf, overrides)
		conf.ImageName = dep

		if _, err := BuildByName(m, conf); err != nil {
			return "", err
		}
	}

	return BuildImage(m, config)
}
