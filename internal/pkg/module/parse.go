package module

import (
	"os"
	"path/filepath"

	"github.com/liamg/peridot/internal/pkg/config"
)

type Module struct {
	Name     string
	Path     string
	Children []Module
	Files    []config.TemplateFile
}

func ParseRoot() (*Module, *config.Config, error) {
	conf, err := config.ParseRoot()
	if err != nil {
		return nil, nil, err
	}
	mod, err := Parse(conf.Path)
	if err != nil {
		return nil, nil, err
	}
	return mod, conf, nil
}

func Parse(path string) (*Module, error) {
	return parse(path, true)
}

func parse(path string, root bool) (*Module, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	conf, err := config.Parse(f)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(filepath.Dir(path))
	if root {
		name = "Root"
	}
	if conf.Name != "" {
		name = conf.Name
	}

	mod := Module{
		Name:  name,
		Path:  filepath.Dir(path),
		Files: conf.TemplateFiles,
	}

	for _, include := range conf.Includes {
		modPath, err := filepath.Rel(mod.Path, include)
		if err != nil {
			return nil, err
		}
		childMod, err := parse(modPath, false)
		if err != nil {
			return nil, err
		}
		mod.Children = append(mod.Children, *childMod)
	}

	return &mod, nil
}
