package builder

import (
	"errors"
	"fmt"
	"strings"

	api "github.com/wabenet/dodo-core/internal/gen-proto/wabenet/dodo/build/v1alpha2"
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
	out := &api.SshAgent{}

	out.SetId(s.ID)
	out.SetIdentityFile(s.IdentityFile)

	return out
}

func ParseSSHAgent(spec string) (SSHAgent, error) {
	out := SSHAgent{}

	switch values := strings.SplitN(spec, "=", 2); len(values) {
	case 2:
		out.ID = values[0]
		out.IdentityFile = values[1]
	default:
		return out, fmt.Errorf("%s: %w", spec, ErrSSHAgentFormat)
	}

	return out, nil
}
