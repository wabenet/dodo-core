package plugin

import (
	"github.com/deckarep/golang-set"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
)

type dependency struct {
	n string
	t string
}

func asDependency(name *api.PluginName) dependency {
	return dependency{n: name.Name, t: name.Type}
}

func ResolveDependencies(pluginMap map[string]map[string]Plugin) ([]Plugin, error) {
	result := []Plugin{}
	names := make(map[dependency]Plugin)
	dependencies := make(map[dependency]mapset.Set)

	for _, ps := range pluginMap {
		for _, p := range ps {
			info := p.PluginInfo()
			if info == nil {
				continue
			}

			deps := mapset.NewSet()
			key := asDependency(info.Name)

			names[key] = p

			for _, dep := range info.Dependencies {
				deps.Add(asDependency(dep))
			}

			dependencies[key] = deps
		}
	}

	for len(dependencies) != 0 {
		withoutDeps := mapset.NewSet()

		for name, deps := range dependencies {
			if deps.Cardinality() == 0 {
				withoutDeps.Add(name)
			}
		}

		if withoutDeps.Cardinality() == 0 {
			return nil, ErrCircularDependency
		}

		for n := range withoutDeps.Iter() {
			delete(dependencies, n.(dependency))
			result = append(result, names[n.(dependency)])
		}

		for n, deps := range dependencies {
			dependencies[n] = deps.Difference(withoutDeps)
		}
	}

	return result, nil
}
