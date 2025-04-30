package builder

import (
	api "github.com/wabenet/dodo-core/api/build/v1alpha2"
)

type BuildSecretsConfig []BuildSecret

type BuildSecret struct {
	ID   string
	Path string
}

func MergeBuildSecretsConfig(first, second BuildSecretsConfig) BuildSecretsConfig {
	return append(first, second...)
}

func BuildSecretsConfigFromProto(b []*api.BuildSecret) BuildSecretsConfig {
	out := BuildSecretsConfig{}

	for _, arg := range b {
		out = append(out, BuildSecretFromProto(arg))
	}

	return out
}

func (b BuildSecretsConfig) ToProto() []*api.BuildSecret {
	out := []*api.BuildSecret{}

	for _, sec := range b {
		out = append(out, sec.ToProto())
	}

	return out
}

func BuildSecretFromProto(b *api.BuildSecret) BuildSecret {
	return BuildSecret{
		ID:   b.GetId(),
		Path: b.GetPath(),
	}
}

func (b BuildSecret) ToProto() *api.BuildSecret {
	return &api.BuildSecret{
		Id:   b.ID,
		Path: b.Path,
	}
}
