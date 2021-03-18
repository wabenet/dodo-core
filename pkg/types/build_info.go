package types

import (
	"reflect"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

func DecodeBuildInfo(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	build := *(target.(**api.BuildInfo))

	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"name":       decoder.String(&build.ImageName),
			"context":    decoder.String(&build.Context),
			"dockerfile": decoder.String(&build.Dockerfile),
			"steps": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(decoder.NewString(), &build.InlineDockerfile),
				reflect.Slice:  decoder.Slice(decoder.NewString(), &build.InlineDockerfile),
			}),
			"inline": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(decoder.NewString(), &build.InlineDockerfile),
				reflect.Slice:  decoder.Slice(decoder.NewString(), &build.InlineDockerfile),
			}),
			"args": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewArgument(), &build.Arguments),
				reflect.Slice:  decoder.Slice(NewArgument(), &build.Arguments),
			}),
			"arguments": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(decoder.NewString(), &build.Arguments),
				reflect.Slice:  decoder.Slice(decoder.NewString(), &build.Arguments),
			}),
			"requires": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(decoder.NewString(), &build.Dependencies),
				reflect.Slice:  decoder.Slice(decoder.NewString(), &build.Dependencies),
			}),
			"dependencies": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(decoder.NewString(), &build.Dependencies),
				reflect.Slice:  decoder.Slice(decoder.NewString(), &build.Dependencies),
			}),
			"secrets": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewSecret(), &build.Secrets),
				reflect.Slice:  decoder.Slice(NewSecret(), &build.Secrets),
			}),
			"ssh": decoder.Kinds(map[reflect.Kind]decoder.Decoding{
				reflect.String: decoder.Singleton(NewSSHAgent(), &build.SshAgents),
				reflect.Slice:  decoder.Slice(NewSSHAgent(), &build.SshAgents),
			}),
			"no_cache":      decoder.Bool(&build.NoCache),
			"force_rebuild": decoder.Bool(&build.ForceRebuild),
			"force_pull":    decoder.Bool(&build.ForcePull),
		}),
	})
}
