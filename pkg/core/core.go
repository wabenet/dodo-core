package core

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
)

func RunByName(overrides *api.Backdrop) error {
	config := AssembleBackdropConfig(overrides.Name, overrides)

	if len(config.ContainerName) == 0 {
		id := make([]byte, 8)
		if _, err := rand.Read(id); err != nil {
			panic(err)
		}

		config.ContainerName = fmt.Sprintf("%s-%s", config.Name, hex.EncodeToString(id))
	}

	if len(config.ImageId) == 0 {
		imageID, err := BuildByName(config.BuildInfo)
		if err != nil {
			return err
		}

		config.ImageId = imageID
	}

	return RunBackdrop(config)
}

func BuildByName(overrides *api.BuildInfo) (string, error) {
	config, err := FindBuildConfig(overrides.ImageName, overrides)
	if err != nil {
		return "", err
	}

	for _, dep := range config.Dependencies {
		if _, err := BuildByName(&api.BuildInfo{ImageName: dep}); err != nil {
			return "", err
		}
	}

	return BuildImage(config)
}
