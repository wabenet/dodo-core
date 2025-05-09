package plugin

import (
	"fmt"
	"strings"

	mapset "github.com/deckarep/golang-set"
	log "github.com/hashicorp/go-hclog"
)

type CircularDependencyError struct {
	Dependencies map[dependency]mapset.Set
}

func (e CircularDependencyError) Error() string {
	lines := []string{"circular dependencies in plugins:"}

	for dep := range e.Dependencies {
		lines = append(lines, fmt.Sprintf("%s/%s", dep.t, dep.n))
	}

	return strings.Join(lines, "\n")
}

type dependency struct {
	n string
	t string
}

func asDependency(id ID) dependency {
	return dependency{n: id.Name, t: id.Type}
}

func ResolveDependencies(pluginMap map[string]map[string]Plugin) []Plugin {
	result := []Plugin{}
	names := make(map[dependency]Plugin)
	dependencies := make(map[dependency]mapset.Set)

	for _, plugins := range pluginMap {
		for _, plugin := range plugins {
			metadata := plugin.Metadata()

			deps := mapset.NewSet()
			key := asDependency(metadata.ID)

			names[key] = plugin

			for _, dep := range metadata.Dependencies {
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
			removed := []string{}
			for dep := range dependencies {
				removed = append(removed, fmt.Sprintf("%s/%s", dep.t, dep.n))
			}

			log.L().Warn("could not fully resolve dependencies, some plugins are automatically removed", "plugins", removed)

			return result
		}

		//nolint:forcetypeassert // we know that the map keys are of type dependency
		for n := range withoutDeps.Iter() {
			delete(dependencies, n.(dependency))
			result = append(result, names[n.(dependency)])
		}

		for n, deps := range dependencies {
			dependencies[n] = deps.Difference(withoutDeps)
		}
	}

	return result
}
