package plugin_test

import (
	"github.com/hashicorp/go-plugin"
	core "github.com/wabenet/dodo-core/api/core/v1alpha5"
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

func (typeA) GRPCClient() (plugin.Plugin, error) {
	return nil, dodo.InvalidError{}
}

func (typeA) GRPCServer(p dodo.Plugin) (plugin.Plugin, error) {
	return nil, dodo.InvalidError{}
}

func (t typeB) String() string {
	return string(t)
}

func (typeB) GRPCClient() (plugin.Plugin, error) {
	return nil, dodo.InvalidError{}
}

func (typeB) GRPCServer(p dodo.Plugin) (plugin.Plugin, error) {
	return nil, dodo.InvalidError{}
}

func (p pluginA) PluginInfo() *core.PluginInfo {
	return &core.PluginInfo{
		Name: &core.PluginName{
			Name: string(p),
			Type: p.Type().String(),
		},
	}
}

func (pluginA) Init() (dodo.Config, error) {
	return map[string]string{}, nil
}

func (pluginA) Cleanup() {}

func (pluginA) Type() dodo.Type {
	return typeAImpl
}

func (p pluginB) PluginInfo() *core.PluginInfo {
	return &core.PluginInfo{
		Name: &core.PluginName{
			Name: string(p),
			Type: p.Type().String(),
		},
		Dependencies: []*core.PluginName{
			pluginAImpl.PluginInfo().Name,
		},
	}
}

func (pluginB) Init() (dodo.Config, error) {
	return map[string]string{}, nil
}

func (pluginB) Cleanup() {}

func (pluginB) Type() dodo.Type {
	return typeBImpl
}

func populatePluginMap() map[string]map[string]dodo.Plugin {
	pluginMap := map[string]map[string]dodo.Plugin{}
	pluginMap[typeAImpl.String()] = map[string]dodo.Plugin{}
	pluginMap[typeBImpl.String()] = map[string]dodo.Plugin{}
	pluginMap[typeAImpl.String()][pluginAImpl.PluginInfo().Name.Name] = pluginAImpl
	pluginMap[typeBImpl.String()][pluginBImpl.PluginInfo().Name.Name] = pluginBImpl

	return pluginMap
}
