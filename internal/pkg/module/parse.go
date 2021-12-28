package module

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/liamg/peridot/internal/pkg/config"
)

func LoadRoot() (Module, error) {
	rootConfig, override, err := config.ParseRoot()
	if err != nil {
		return nil, err
	}
	rootConfig.Name = "root"
	return loadModule(*rootConfig, config.BaseVariables(), override)
}

func loadModule(conf config.Module, combined map[string]interface{}, override *config.Override) (Module, error) {

	if conf.Name == "" {
		conf.Name = filepath.Base(filepath.Dir(conf.Dir))
	}

	mod := module{
		conf:      conf,
		variables: combined,
	}

	// load and set mod.files
	for _, f := range conf.Files {
		file, err := loadFile(conf, f, combined)
		if err != nil {
			return nil, err
		}
		mod.files = append(mod.files, file)
	}

	for _, childInfo := range conf.Modules {
		child, err := loadModuleFromSource(childInfo, conf, override)
		if err != nil {
			return nil, err
		}
		mod.children = append(mod.children, child)
	}

	if err := mod.Validate(); err != nil {
		return nil, err
	}

	return &mod, nil
}

func mergeVars(defaults []config.Variable, base, configured, overrides map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for _, def := range defaults {
		if !def.Required {
			merged[def.Name] = def.Default
		}
	}
	for _, m := range []map[string]interface{}{base, configured, overrides} {
		for key, val := range m {
			merged[key] = val
		}
	}
	return merged
}

func loadModuleFromSource(info config.InnerModule, parent config.Module, override *config.Override) (Module, error) {

	var path string

	switch {
	case strings.HasPrefix(info.Source, "./"):
		var err error
		path, err = filepath.Abs(filepath.Join(parent.Dir, info.Source))
		if err != nil {
			return nil, err
		}
		path = filepath.Join(path, config.Filename)
	default:
		return nil, fmt.Errorf("invalid module source '%s' - local modules should begin with './'. "+
			"To load external modules, use a full URL, or to use built in modules, use 'builtin:NAME'", info.Source)
	}

	variableDefaults, err := config.ParseVariables(path)
	if err != nil {
		return nil, err
	}

	combined := mergeVars(variableDefaults, config.BaseVariables(), info.Variables, override.Variables[info.Name])

	conf, err := config.Parse(path, combined)
	if err != nil {
		return nil, err
	}
	conf.Name = info.Name
	return loadModule(*conf, combined, override)
}
