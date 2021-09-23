package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/configuration"
)

func RunByName(m plugin.Manager, overrides *api.Backdrop) error {
	config := configuration.AssembleBackdropConfig(m, overrides.Name, overrides)

	if len(config.ContainerName) == 0 {
		id := make([]byte, 8)
		if _, err := rand.Read(id); err != nil {
			panic(err)
		}

		config.ContainerName = fmt.Sprintf("%s-%s", config.Name, hex.EncodeToString(id))
	}

	if len(config.ImageId) == 0 {
		for _, dep := range config.BuildInfo.Dependencies {
			if _, err := BuildByName(m, &api.BuildInfo{ImageName: dep}); err != nil {
				return err
			}
		}

		imageID, err := BuildImage(m, config.BuildInfo)
		if err != nil {
			return err
		}

		config.ImageId = imageID
	}

	return RunBackdrop(m, config)
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
