package configuration

import (
	"fmt"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	core "github.com/wabenet/dodo-core/api/core/v1alpha5"
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

func AssembleBackdropConfig(m plugin.Manager, name string, overrides *core.Backdrop) (*core.Backdrop, error) {
	var errs error

	config := &core.Backdrop{Entrypoint: &core.Entrypoint{}}
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

func FindBuildConfig(m plugin.Manager, name string, overrides *core.BuildInfo) (*core.BuildInfo, error) {
	var errs error

	for n, p := range m.GetPlugins(Type.String()) {
		log.L().Debug("Fetching configuration from plugin", "name", n)

		confs, err := p.(Configuration).ListBackdrops()
		if err != nil {
			errs = multierror.Append(errs, err)

			continue
		}

		for _, conf := range confs {
			if conf.GetBuildInfo() != nil && conf.GetBuildInfo().GetImageName() == name {
				config := &core.BuildInfo{}
				MergeBuildInfo(config, conf.GetBuildInfo())
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

func MergeBackdrop(target, source *core.Backdrop) {
	if len(source.GetName()) > 0 {
		target.Name = source.GetName()
	}

	target.Aliases = append(target.GetAliases(), source.GetAliases()...)

	if len(source.GetImageId()) > 0 {
		target.ImageId = source.GetImageId()
	}

	if len(source.GetRuntime()) > 0 {
		target.Runtime = source.GetRuntime()
	}

	if source.GetEntrypoint() != nil {
		if source.GetEntrypoint().GetInteractive() {
			target.Entrypoint.Interactive = true
		}

		if len(source.GetEntrypoint().GetInterpreter()) > 0 {
			target.Entrypoint.Interpreter = source.GetEntrypoint().GetInterpreter()
		}

		if len(source.GetEntrypoint().GetScript()) > 0 {
			target.Entrypoint.Script = source.GetEntrypoint().GetScript()
		}

		if len(source.GetEntrypoint().GetArguments()) > 0 {
			target.Entrypoint.Arguments = source.GetEntrypoint().GetArguments()
		}
	}

	if len(source.GetContainerName()) > 0 {
		target.ContainerName = source.GetContainerName()
	}

	target.Environment = append(target.GetEnvironment(), source.GetEnvironment()...)

	if len(source.GetUser()) > 0 {
		target.User = source.GetUser()
	}

	target.Volumes = append(target.GetVolumes(), source.GetVolumes()...)
	target.Devices = append(target.GetDevices(), source.GetDevices()...)
	target.Ports = append(target.GetPorts(), source.GetPorts()...)
	target.Capabilities = append(target.GetCapabilities(), source.GetCapabilities()...)

	if len(source.GetWorkingDir()) > 0 {
		target.WorkingDir = source.GetWorkingDir()
	}

	if source.GetBuildInfo() != nil {
		if target.GetBuildInfo() == nil {
			target.BuildInfo = source.GetBuildInfo()
		} else {
			MergeBuildInfo(target.GetBuildInfo(), source.GetBuildInfo())
		}
	}
}

func MergeBuildInfo(target, source *core.BuildInfo) {
	if len(source.GetBuilder()) > 0 {
		target.Builder = source.GetBuilder()
	}

	if len(source.GetImageName()) > 0 {
		target.ImageName = source.GetImageName()
	}

	if len(source.GetContext()) > 0 {
		target.Context = source.GetContext()
	}

	if len(source.GetDockerfile()) > 0 {
		target.Dockerfile = source.GetDockerfile()
	}

	if len(source.GetInlineDockerfile()) > 0 {
		target.InlineDockerfile = source.GetInlineDockerfile()
	}

	target.Arguments = append(target.GetArguments(), source.GetArguments()...)
	target.Secrets = append(target.GetSecrets(), source.GetSecrets()...)
	target.SshAgents = append(target.GetSshAgents(), source.GetSshAgents()...)

	if source.GetNoCache() {
		target.NoCache = true
	}

	if source.GetForceRebuild() {
		target.ForceRebuild = true
	}

	if source.GetForcePull() {
		target.ForcePull = true
	}

	target.Dependencies = append(target.GetDependencies(), source.GetDependencies()...)
}

func ValidateBackdrop(b *core.Backdrop) error {
	if b.GetImageId() == "" && b.GetBuildInfo() == nil {
		return InvalidError{Name: b.GetName(), Message: "neither image nor build configured"}
	}

	if b.GetBuildInfo() != nil {
		if err := ValidateBuildInfo(b.GetBuildInfo()); err != nil {
			return err
		}
	}

	return nil
}

func ValidateBuildInfo(b *core.BuildInfo) error {
	return nil
}
