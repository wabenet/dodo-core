package core

import (
	"fmt"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/builder"
	"github.com/dodo-cli/dodo-core/pkg/plugin/configuration"
	"github.com/dodo-cli/dodo-core/pkg/plugin/runtime"
	log "github.com/hashicorp/go-hclog"
)

func GetRuntime(name string) (runtime.ContainerRuntime, error) {
	for n, p := range plugin.GetPlugins(runtime.Type.String()) {
		if name != "" && name != n {
			continue
		}

		if rt, ok := p.(runtime.ContainerRuntime); ok {
			return rt, nil
		}
	}

	return nil, fmt.Errorf("could not find container runtime: %w", plugin.ErrNoValidPluginFound)
}

func GetBuilder(name string) (builder.ImageBuilder, error) {
	for n, p := range plugin.GetPlugins(builder.Type.String()) {
		if name != "" && name != n {
			continue
		}

		if rt, ok := p.(builder.ImageBuilder); ok {
			return rt, nil
		}
	}

	return nil, fmt.Errorf("could not find image builder: %w", plugin.ErrNoValidPluginFound)
}

func AssembleBackdropConfig(name string, overrides *api.Backdrop) *api.Backdrop {
	config := &api.Backdrop{Entrypoint: &api.Entrypoint{}}

	for _, p := range plugin.GetPlugins(configuration.Type.String()) {
		info, err := p.PluginInfo()
		if err != nil {
			log.L().Warn("could not read plugin info")
			continue
		}

		log.L().Debug("Fetching configuration from plugin", "name", info.Name)

		conf, err := p.(configuration.Configuration).GetBackdrop(name)
		if err != nil {
			log.L().Warn("could not get config", "error", err)
			continue
		}

		mergeBackdrop(config, conf)
	}

	mergeBackdrop(config, overrides)
	log.L().Debug("assembled configuration", "backdrop", config)

	return config
}

func FindBuildConfig(name string, overrides *api.BuildInfo) (*api.BuildInfo, error) {
	for _, p := range plugin.GetPlugins(configuration.Type.String()) {
		info, err := p.PluginInfo()
		if err != nil {
			log.L().Warn("could not read plugin info")
			continue
		}

		log.L().Debug("Fetching configuration from plugin", "name", info.Name)

		conf, err := p.(configuration.Configuration).GetBackdrop(name)
		if err != nil {
			log.L().Warn("could not get config", "error", err)
			continue
		}

		if conf.BuildInfo != nil && conf.BuildInfo.ImageName == name {
			return conf.BuildInfo, nil
		}
	}

	return nil, fmt.Errorf("no config found for image")
}

func mergeBackdrop(target *api.Backdrop, source *api.Backdrop) {
	if len(source.Name) > 0 {
		target.Name = source.Name
	}

	target.Aliases = append(target.Aliases, source.Aliases...)

	if len(source.ImageId) > 0 {
		target.ImageId = source.ImageId
	}

	if len(source.Runtime) > 0 {
		target.Runtime = source.Runtime
	}

	if source.Entrypoint != nil {
		if source.Entrypoint.Interactive {
			target.Entrypoint.Interactive = true
		}

		if len(source.Entrypoint.Interpreter) > 0 {
			target.Entrypoint.Interpreter = source.Entrypoint.Interpreter
		}

		if len(source.Entrypoint.Script) > 0 {
			target.Entrypoint.Script = source.Entrypoint.Script
		}

		if len(source.Entrypoint.Arguments) > 0 {
			target.Entrypoint.Arguments = source.Entrypoint.Arguments
		}
	}

	if len(source.ContainerName) > 0 {
		target.ContainerName = source.ContainerName
	}

	target.Environment = append(target.Environment, source.Environment...)

	if len(source.User) > 0 {
		target.User = source.User
	}

	target.Volumes = append(target.Volumes, source.Volumes...)
	target.Devices = append(target.Devices, source.Devices...)
	target.Ports = append(target.Ports, source.Ports...)
	target.Capabilities = append(target.Capabilities, source.Capabilities...)

	if len(source.WorkingDir) > 0 {
		target.WorkingDir = source.WorkingDir
	}

	if source.BuildInfo != nil {
		if target.BuildInfo == nil {
			target.BuildInfo = source.BuildInfo
		} else {
			mergeBuildInfo(target.BuildInfo, source.BuildInfo)
		}
	}
}

func mergeBuildInfo(target *api.BuildInfo, source *api.BuildInfo) {
	if len(source.Builder) > 0 {
		target.Builder = source.Builder
	}

	if len(source.ImageName) > 0 {
		target.ImageName = source.ImageName
	}

	if len(source.Context) > 0 {
		target.Context = source.Context
	}

	if len(source.Dockerfile) > 0 {
		target.Dockerfile = source.Dockerfile
	}

	if len(source.InlineDockerfile) > 0 {
		target.InlineDockerfile = source.InlineDockerfile
	}

	target.Arguments = append(target.Arguments, source.Arguments...)
	target.Secrets = append(target.Secrets, source.Secrets...)
	target.SshAgents = append(target.SshAgents, source.SshAgents...)

	if source.NoCache {
		target.NoCache = true
	}

	if source.ForceRebuild {
		target.ForceRebuild = true
	}

	if source.ForcePull {
		target.ForcePull = true
	}

	target.Dependencies = append(target.Dependencies, source.Dependencies...)
}
