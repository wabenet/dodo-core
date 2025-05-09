package plugin

import (
	"io"

	"github.com/hashicorp/go-plugin"
	api "github.com/wabenet/dodo-core/api/plugin/v1alpha2"
)

type Plugin interface {
	Type() Type
	Metadata() Metadata
	Init() (Config, error)
	Cleanup()
}

type Type interface {
	String() string
	GRPCClient() (plugin.Plugin, error)
	GRPCServer(p Plugin) (plugin.Plugin, error)
}

type Metadata struct {
	ID           ID
	Dependencies Dependencies
	Labels       Labels
}

func NewMetadata(pluginType Type, name string) Metadata {
	return Metadata{
		ID: ID{
			Type: pluginType.String(),
			Name: name,
		},
		Dependencies: Dependencies{},
		Labels:       Labels{},
	}
}

func (m Metadata) WithDependencies(deps ...Plugin) Metadata {
	for _, plugin := range deps {
		m.Dependencies = append(m.Dependencies, plugin.Metadata().ID)
	}

	return m
}

func (m Metadata) WithLabels(labels Labels) Metadata {
	m.Labels = labels

	return m
}

func MetadataFromProto(m *api.PluginMetadata) Metadata {
	return Metadata{
		ID:           IDFromProto(m.GetId()),
		Dependencies: DependenciesFromProto(m.GetDependencies()),
		Labels:       m.GetLabels(),
	}
}

func (m Metadata) ToProto() *api.PluginMetadata {
	out := &api.PluginMetadata{}

	out.SetId(m.ID.ToProto())
	out.SetDependencies(m.Dependencies.ToProto())
	out.SetLabels(m.Labels)

	return out
}

type ID struct {
	Type string
	Name string
}

func IDFromProto(id *api.PluginID) ID {
	return ID{
		Type: id.GetType(),
		Name: id.GetName(),
	}
}

func (id ID) ToProto() *api.PluginID {
	out := &api.PluginID{}

	out.SetType(id.Type)
	out.SetName(id.Name)

	return out
}

type Dependencies []ID

func DependenciesFromProto(deps []*api.PluginID) Dependencies {
	out := Dependencies{}

	for _, dep := range deps {
		out = append(out, IDFromProto(dep))
	}

	return out
}

func (d Dependencies) ToProto() []*api.PluginID {
	out := []*api.PluginID{}

	for _, dep := range d {
		out = append(out, dep.ToProto())
	}

	return out
}

type Labels map[string]string

type Config map[string]string

type StreamConfig struct {
	Stdin          io.Reader
	Stdout         io.Writer
	Stderr         io.Writer
	TerminalHeight uint32
	TerminalWidth  uint32
}
