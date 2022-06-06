package configuration

import (
	"fmt"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	api "github.com/wabenet/dodo-core/api/v1alpha3"
	"github.com/wabenet/dodo-core/pkg/plugin"
)

type NotFoundError struct {
	Name   string
	Reason error
}

func (e NotFoundError) Error() string {
	if e.Reason == nil {
		return fmt.Sprintf(
			"could not find any configuration for '%s'", e.Name,
		)
	}

	return fmt.Sprintf(
		"could not find any configuration for '%s': %s",
		e.Name,
		e.Reason.Error(),
	)
}

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

func AssembleBackdropConfig(m plugin.Manager, name string, overrides *api.Backdrop) (*api.Backdrop, error) {
	var errs error

	config := &api.Backdrop{Entrypoint: &api.Entrypoint{}}
	foundSomething := false

	for n, p := range m.GetPlugins(Type.String()) {
		log.L().Debug("Fetching configuration from plugin", "name", n)

		conf, err := p.(Configuration).GetBackdrop(name)
		if err != nil {
			errs = multierror.Append(errs, err)

			continue
		}

		foundSomething = true

		MergeBackdrop(config, conf)
	}

	if !foundSomething {
		return nil, NotFoundError{Name: name, Reason: errs}
	}

	MergeBackdrop(config, overrides)
	log.L().Debug("assembled configuration", "backdrop", config)

	err := ValidateBackdrop(config)

	return config, err
}

func FindBuildConfig(m plugin.Manager, name string, overrides *api.BuildInfo) (*api.BuildInfo, error) {
	var errs error

	for n, p := range m.GetPlugins(Type.String()) {
		log.L().Debug("Fetching configuration from plugin", "name", n)

		confs, err := p.(Configuration).ListBackdrops()
		if err != nil {
			errs = multierror.Append(errs, err)

			continue
		}

		for _, conf := range confs {
			if conf.BuildInfo != nil && conf.BuildInfo.ImageName == name {
				config := &api.BuildInfo{}
				MergeBuildInfo(config, conf.BuildInfo)
				MergeBuildInfo(config, overrides)

				if err := ValidateBuildInfo(config); err != nil {
					errs = multierror.Append(errs, err)

					continue
				}

				return config, nil
			}
		}
	}

	return nil, NotFoundError{Name: name, Reason: errs}
}

func MergeBackdrop(target *api.Backdrop, source *api.Backdrop) {
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
			MergeBuildInfo(target.BuildInfo, source.BuildInfo)
		}
	}
}

func MergeBuildInfo(target *api.BuildInfo, source *api.BuildInfo) {
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

func ValidateBackdrop(b *api.Backdrop) error {
	if b.ImageId == "" && b.BuildInfo == nil {
		return InvalidError{Name: b.Name, Message: "neither image nor build configured"}
	}

	if b.BuildInfo != nil {
		if err := ValidateBuildInfo(b.BuildInfo); err != nil {
			return err
		}
	}

	return nil
}

func ValidateBuildInfo(b *api.BuildInfo) error {
	return nil
}
