package builder

import (
	api "github.com/wabenet/dodo-core/api/build/v1alpha2"
)

type BuildConfig struct {
	ImageName        string
	Context          string
	Dockerfile       string
	InlineDockerfile []string

	Arguments BuildArgumentsConfig
	Secrets   BuildSecretsConfig
	SSHAgents SSHAgentConfig

	NoCache      bool
	ForceRebuild bool
	ForcePull    bool

	Dependencies []string
}

func EmptyBuildConfig() BuildConfig {
	return BuildConfig{
		Arguments: BuildArgumentsConfig{},
		Secrets:   BuildSecretsConfig{},
		SSHAgents: SSHAgentConfig{},
	}
}

func MergeBuildConfig(first, second BuildConfig) BuildConfig {
	result := first

	if len(second.ImageName) > 0 {
		result.ImageName = second.ImageName
	}

	if len(second.Context) > 0 {
		result.Context = second.Context
	}

	if len(second.Dockerfile) > 0 {
		result.Dockerfile = second.Dockerfile
	}

	if len(second.InlineDockerfile) > 0 {
		result.InlineDockerfile = second.InlineDockerfile
	}

	result.Arguments = MergeBuildArgumentsConfig(first.Arguments, second.Arguments)
	result.Secrets = MergeBuildSecretsConfig(first.Secrets, second.Secrets)
	result.SSHAgents = MergeSSHAgentConfig(first.SSHAgents, second.SSHAgents)

	if second.NoCache {
		result.NoCache = true
	}

	if second.ForceRebuild {
		result.ForceRebuild = true
	}

	if second.ForcePull {
		result.ForcePull = true
	}

	result.Dependencies = append(first.Dependencies, second.Dependencies...)

	return result
}

func BuildConfigFromProto(b *api.BuildConfig) BuildConfig {
	return BuildConfig{
		ImageName:        b.GetImageName(),
		Context:          b.GetContext(),
		Dockerfile:       b.GetDockerfile(),
		InlineDockerfile: b.GetInlineDockerfile(),
		Arguments:        BuildArgumentsConfigFromProto(b.GetArguments()),
		Secrets:          BuildSecretsConfigFromProto(b.GetSecrets()),
		SSHAgents:        SSHAgentConfigFromProto(b.GetSshAgents()),
		NoCache:          b.GetNoCache(),
		ForceRebuild:     b.GetForceRebuild(),
		ForcePull:        b.GetForcePull(),
		Dependencies:     b.GetDependencies(),
	}
}

func (b BuildConfig) ToProto() *api.BuildConfig {
	out := &api.BuildConfig{}

	out.SetImageName(b.ImageName)
	out.SetContext(b.Context)
	out.SetDockerfile(b.Dockerfile)
	out.SetInlineDockerfile(b.InlineDockerfile)
	out.SetArguments(b.Arguments.ToProto())
	out.SetSecrets(b.Secrets.ToProto())
	out.SetSshAgents(b.SSHAgents.ToProto())
	out.SetNoCache(b.NoCache)
	out.SetForceRebuild(b.ForceRebuild)
	out.SetForcePull(b.ForcePull)
	out.SetDependencies(b.Dependencies)

	return out
}
