package plugin_test

import (
	"github.com/hashicorp/go-plugin"
	dodo "github.com/wabenet/dodo-core/pkg/plugin"
)

type (
	typeA   string
	typeB   string
	pluginA string
	pluginB string
)

const (
	typeAImpl   typeA   = "typeA"
	typeBImpl   typeB   = "tybeB"
	pluginAImpl pluginA = "pluginA"
	pluginBImpl pluginB = "pluginB"
)

func (t typeA) String() string {
	return string(t)
}

func (typeA) GRPCClient() (plugin.Plugin, error) { //nolint:ireturn
	return nil, dodo.InvalidError{}
}

func (typeA) GRPCServer(_ dodo.Plugin) (plugin.Plugin, error) { //nolint:ireturn
	return nil, dodo.InvalidError{}
}

func (t typeB) String() string {
	return string(t)
}

func (typeB) GRPCClient() (plugin.Plugin, error) { //nolint:ireturn
	return nil, dodo.InvalidError{}
}

func (typeB) GRPCServer(_ dodo.Plugin) (plugin.Plugin, error) { //nolint:ireturn
	return nil, dodo.InvalidError{}
}

func (p pluginA) Metadata() dodo.Metadata {
	return dodo.NewMetadata(p.Type(), "")
}

func (pluginA) Init() (dodo.Config, error) {
	return map[string]string{}, nil
}

func (pluginA) Cleanup() {}

func (pluginA) Type() dodo.Type { //nolint:ireturn
	return typeAImpl
}

func (p pluginB) Metadata() dodo.Metadata {
	return dodo.NewMetadata(p.Type(), "").WithDependencies(pluginAImpl)
}

func (pluginB) Init() (dodo.Config, error) {
	return map[string]string{}, nil
}

func (pluginB) Cleanup() {}

func (pluginB) Type() dodo.Type { //nolint:ireturn
	return typeBImpl
}

func populatePluginMap() map[string]map[string]dodo.Plugin {
	pluginMap := map[string]map[string]dodo.Plugin{}
	pluginMap[typeAImpl.String()] = map[string]dodo.Plugin{}
	pluginMap[typeBImpl.String()] = map[string]dodo.Plugin{}
	pluginMap[typeAImpl.String()][pluginAImpl.Metadata().ID.Name] = pluginAImpl
	pluginMap[typeBImpl.String()][pluginBImpl.Metadata().ID.Name] = pluginBImpl

	return pluginMap
}
