package types

import (
	"fmt"
	"reflect"
	"strings"

	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
)

const ErrSSHAgentFormat FormatError = "invalid ssh agent format"

func ParseSSHAgent(spec string) (*api.SshAgent, error) {
	agent := &api.SshAgent{}

	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 2:
		agent.Id = values[0]
		agent.IdentityFile = values[1]
	default:
		return nil, fmt.Errorf("%s: %w", spec, ErrSSHAgentFormat)
	}

	return agent, nil
}

func NewSSHAgent() decoder.Producer {
	return func() (interface{}, decoder.Decoding) {
		target := &api.SshAgent{}
		return &target, DecodeSSHAgent(&target)
	}
}

func DecodeSSHAgent(target interface{}) decoder.Decoding {
	// TODO: wtf this cast
	agent := *(target.(**api.SshAgent))

	return decoder.Kinds(map[reflect.Kind]decoder.Decoding{
		reflect.Map: decoder.Keys(map[string]decoder.Decoding{
			"id":   decoder.String(&agent.Id),
			"file": decoder.String(&agent.IdentityFile),
		}),
		reflect.String: func(d *decoder.Decoder, config interface{}) {
			var decoded string

			decoder.String(&decoded)(d, config)

			if a, err := ParseSSHAgent(decoded); err != nil {
				d.Error(err)
			} else {
				agent.Id = a.Id
				agent.IdentityFile = a.IdentityFile
			}
		},
	})
}
