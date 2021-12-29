package module

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/variable"
)

func LoadRoot() (Module, error) {
	rootConfig, override, err := config.ParseRoot()
	if err != nil {
		return nil, err
	}
	rootConfig.Name = "root"
	return loadModule(*rootConfig, config.BaseVariables(), override)
}

func loadModule(conf config.Module, combined variable.Collection, override *config.Override) (Module, error) {

	if conf.Name == "" {
		conf.Name = filepath.Base(conf.Dir)
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
		if !filtersMatch(childInfo.Filters) {
			continue
		}
		child, err := loadModuleFromSource(childInfo, conf, override)
		if err != nil {
			return nil, err
		}
		mod.children = append(mod.children, child)
	}

	if err := mod.Validate(); err != nil {
		return nil, fmt.Errorf("validation error in module '%s': %w", mod.conf.Name, err)
	}
	return &mod, nil
}

func mergeVars(defaults []config.Variable, base, configured, overrides variable.Collection) variable.Collection {
	merged := variable.NewCollection(nil)
	for _, def := range defaults {
		if !def.Required {
			merged.Set(def.Name, def.Default)
		}
	}
	for _, m := range []variable.Collection{base, configured, overrides} {
		for key, val := range m.AsMap() {
			merged.Set(key, val)
		}
	}
	return merged
}

func loadModuleFromSource(info config.InnerModule, parent config.Module, override *config.Override) (Module, error) {

	var path string

	switch {
	case strings.HasPrefix(info.Source, "builtin:"): // builtin modules
		combined := mergeVars(
			nil,
			config.BaseVariables(),
			variable.NewCollection(info.Variables),
			variable.NewCollection(override.Variables).Get(info.Name).AsCollection(),
		)
		return loadBuiltin(info.Source[8:], info.Name, combined)
	case strings.HasPrefix(info.Source, "./"): // locally defined modules
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

	combined := mergeVars(
		variableDefaults,
		config.BaseVariables(),
		variable.NewCollection(info.Variables),
		variable.NewCollection(override.Variables).Get(info.Name).AsCollection(),
	)

	conf, err := config.Parse(path, combined)
	if err != nil {
		return nil, err
	}
	conf.Name = info.Name
	return loadModule(*conf, combined, override)
}
