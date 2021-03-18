package core

import (
	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
)

func BuildImage(config *api.BuildInfo) (string, error) {
	b, err := GetBuilder(config.Builder)
	if err != nil {
		return "", err
	}

	return b.CreateImage(config)
}
