package module

import (
	"fmt"
	"strings"

	"github.com/liamg/peridot/internal/pkg/config"
)

func sortDependencies(modules []config.InnerModule) ([]config.InnerModule, error) {
	var sorted []config.InnerModule

	handled := make([]bool, len(modules))

	// process all modules without dependences first
	for i, m := range modules {
		if len(m.DependsOn) == 0 {
			sorted = append(sorted, m)
			handled[i] = true
			continue
		}
	}

	// process all modules with dependencies satisfied repeatedly until we are done or can go no further
	for len(sorted) < len(modules) {
		last := len(sorted)
		for i, m := range modules {
			if handled[i] {
				continue
			}
			missingDeps := findUnmetDependencies(m, sorted)
			if len(missingDeps) == 0 {
				sorted = append(sorted, m)
				handled[i] = true
			}
		}
		if last == len(sorted) {
			if len(sorted) == len(modules) {
				break
			}
			for i, val := range handled {
				if !val {
					missingDeps := findUnmetDependencies(modules[i], sorted)
					for _, missing := range missingDeps {
						var exists bool
						for _, each := range modules {
							if each.Name == missing {
								exists = true
								break
							}
						}
						if !exists {
							return nil, fmt.Errorf("failed to satisfy dependencies for module '%s': module '%s' was requested but does not exist", modules[i].Name, missing)
						}

					}
					return nil, fmt.Errorf("failed to satisfy dependencies for module '%s' (needs unmet: %s)", modules[i].Name, strings.Join(missingDeps, ", "))
				}
			}
			// should be impossible...
			return nil, fmt.Errorf("bad dependency chain")
		}
	}

	return sorted, nil
}

func findUnmetDependencies(m config.InnerModule, met []config.InnerModule) []string {
	var unmet []string
	for _, dep := range m.DependsOn {
		var found bool
		for _, done := range met {
			if done.Name == dep {
				found = true
				break
			}
		}
		if !found {
			unmet = append(unmet, dep)
		}
	}
	return unmet
}
