package configuration

import (
	"fmt"

	api "github.com/wabenet/dodo-core/api/configuration/v1alpha2"
	"github.com/wabenet/dodo-core/pkg/plugin/builder"
	"github.com/wabenet/dodo-core/pkg/plugin/runtime"
)

type InvalidError struct {
	Name    string
	Message string
}

func (e InvalidError) Error() string {
	return fmt.Sprintf(
		"invalid configuration for '%s': %s",
		e.Name,
		e.Message,
	)
}

type Backdrop struct {
	Name    string
	Aliases []string
	Runtime string
	Builder string

	ContainerConfig runtime.ContainerConfig
	BuildConfig     builder.BuildConfig

	RequiredFiles FilesConfig
}

func EmptyBackdrop() Backdrop {
	return Backdrop{
		ContainerConfig: runtime.EmptyContainerConfig(),
		BuildConfig:     builder.EmptyBuildConfig(),
	}
}

func MergeBackdrops(first, second Backdrop) Backdrop {
	result := first

	if len(second.Name) > 0 {
		result.Name = second.Name
	}

	result.Aliases = append(first.Aliases, second.Aliases...)

	if len(second.Runtime) > 0 {
		result.Runtime = second.Runtime
	}

	if len(second.Builder) > 0 {
		result.Builder = second.Builder
	}

	result.ContainerConfig = runtime.MergeContainerConfig(first.ContainerConfig, second.ContainerConfig)
	result.BuildConfig = builder.MergeBuildConfig(first.BuildConfig, second.BuildConfig)
	result.RequiredFiles = MergeFilesConfig(first.RequiredFiles, second.RequiredFiles)

	return result
}

func (b Backdrop) Validate() error {
	if len(b.ContainerConfig.Image)+len(b.BuildConfig.Dockerfile)+len(b.BuildConfig.InlineDockerfile) == 0 {
		return InvalidError{Name: b.Name, Message: "neither image nor build configured"}
	}

	return nil
}

func BackdropFromProto(b *api.Backdrop) Backdrop {
	return Backdrop{
		Name:            b.GetName(),
		Aliases:         b.GetAliases(),
		Runtime:         b.GetRuntime(),
		Builder:         b.GetBuilder(),
		ContainerConfig: runtime.ContainerConfigFromProto(b.GetContainerConfig()),
		BuildConfig:     builder.BuildConfigFromProto(b.GetBuildConfig()),
		RequiredFiles:   FilesConfigFromProto(b.GetRequiredFiles()),
	}
}

func (b Backdrop) ToProto() *api.Backdrop {
	return &api.Backdrop{
		Name:            b.Name,
		Aliases:         b.Aliases,
		Runtime:         b.Runtime,
		Builder:         b.Builder,
		ContainerConfig: b.ContainerConfig.ToProto(),
		BuildConfig:     b.BuildConfig.ToProto(),
		RequiredFiles:   b.RequiredFiles.ToProto(),
	}
}
