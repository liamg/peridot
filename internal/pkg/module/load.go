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

func loadModuleFromSource(info config.InnerModule, parent config.Module, override *config.Override) (Module, error) {

	var path string

	combined := variable.NewCollection(info.Variables)
	if override != nil {
		overrides := variable.NewCollection(override.Variables).Get(info.Name).AsCollection()
		combined.MergeIn(overrides)
	}

	switch {
	case strings.HasPrefix(info.Source, "builtin:"): // builtin modules
		return loadBuiltin(info.Source[8:], info.Name, combined)
	case strings.HasPrefix(info.Source, "./"): // locally defined modules
		var err error
		path, err = filepath.Abs(filepath.Join(parent.Dir, info.Source))
		if err != nil {
			return nil, err
		}
		path = filepath.Join(path, config.Filename)
		return loadModuleFromPath(path, info.Name, combined)
	default:
		return nil, fmt.Errorf("invalid module source '%s' - local modules should begin with './'. "+
			"To load external modules, use a full URL, or to use built in modules, use 'builtin:NAME'", info.Source)
	}
}

func loadModuleFromPath(path string, name string, combined variable.Collection) (Module, error) {
	variableDefaults, err := config.ParseVariables(path)
	if err != nil {
		return nil, err
	}

	combined = applyVariableDefaults(variableDefaults, combined)

	conf, err := config.Parse(path, combined)
	if err != nil {
		return nil, err
	}
	conf.Name = name
	return loadModule(*conf, combined, nil)
}
