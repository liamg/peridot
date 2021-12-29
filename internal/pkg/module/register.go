package module

import (
	"fmt"
	"sync"

	"github.com/liamg/peridot/internal/pkg/variable"
)

var moduleRegistry = struct {
	sync.Mutex
	modules map[string]BuiltIn
}{
	modules: make(map[string]BuiltIn),
}

type BuiltIn interface {
	Module
	ApplyVariables(vars variable.Collection) error
}

func RegisterBuiltin(name string, builtin BuiltIn) {
	moduleRegistry.Lock()
	defer moduleRegistry.Unlock()
	if _, exists := moduleRegistry.modules[name]; exists {
		panic(fmt.Sprintf("cannot register multiple builtin modules with the same name: '%s'", name))
	}
	moduleRegistry.modules[name] = builtin
}

func loadBuiltin(builtin, name string, vars variable.Collection) (Module, error) {
	moduleRegistry.Lock()
	defer moduleRegistry.Unlock()
	if m, exists := moduleRegistry.modules[builtin]; exists {
		if err := m.ApplyVariables(vars); err != nil {
			return nil, err
		}
		if err := m.Validate(); err != nil {
			return nil, err
		}
		return m, nil
	}
	return nil, fmt.Errorf("builtin module does not exist: '%s'", name)
}
