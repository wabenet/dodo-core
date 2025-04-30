package builder

import (
	"errors"
	"fmt"
	"strings"

	api "github.com/wabenet/dodo-core/api/build/v1alpha2"
)

var ErrSSHAgentFormat = errors.New("invalid ssh agent format")

type SSHAgentConfig []SSHAgent

type SSHAgent struct {
	ID           string
	IdentityFile string
}

func MergeSSHAgentConfig(first, second SSHAgentConfig) SSHAgentConfig {
	return append(first, second...)
}

func SSHAgentConfigFromProto(b []*api.SshAgent) SSHAgentConfig {
	out := SSHAgentConfig{}

	for _, arg := range b {
		out = append(out, SSHAgentFromProto(arg))
	}

	return out
}

func (s SSHAgentConfig) ToProto() []*api.SshAgent {
	out := []*api.SshAgent{}

	for _, ssh := range s {
		out = append(out, ssh.ToProto())
	}

	return out
}

func SSHAgentFromProto(s *api.SshAgent) SSHAgent {
	return SSHAgent{
		ID:           s.GetId(),
		IdentityFile: s.GetIdentityFile(),
	}
}

func (s SSHAgent) ToProto() *api.SshAgent {
	return &api.SshAgent{
		Id:           s.ID,
		IdentityFile: s.IdentityFile,
	}
}

func ParseSSHAgent(spec string) (SSHAgent, error) {
	agent := SSHAgent{}

	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 2:
		agent.ID = values[0]
		agent.IdentityFile = values[1]
	default:
		return agent, fmt.Errorf("%s: %w", spec, ErrSSHAgentFormat)
	}

	return agent, nil
}
