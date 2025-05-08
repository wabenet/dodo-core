package plugin_test

import (
	"github.com/hashicorp/go-plugin"
	api "github.com/wabenet/dodo-core/api/plugin/v1alpha1"
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

func (p pluginA) PluginInfo() *api.PluginInfo {
	return dodo.MkInfo(p.Type(), "")
}

func (pluginA) Init() (dodo.Config, error) {
	return map[string]string{}, nil
}

func (pluginA) Cleanup() {}

func (pluginA) Type() dodo.Type {
	return typeAImpl
}

func (p pluginB) PluginInfo() *api.PluginInfo {
	info := dodo.MkInfo(p.Type(), "")

	info.SetDependencies([]*api.PluginName{
		pluginAImpl.PluginInfo().GetName(),
	})

	return info
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
	pluginMap[typeAImpl.String()][pluginAImpl.PluginInfo().GetName().GetName()] = pluginAImpl
	pluginMap[typeBImpl.String()][pluginBImpl.PluginInfo().GetName().GetName()] = pluginBImpl

	return pluginMap
}
