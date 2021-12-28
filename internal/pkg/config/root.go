package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Dir   string `yaml:"-"`
	Path  string `yaml:"-"`
	Debug bool   `yaml:"debug"`
}

type Override struct {
	Variables map[string]map[string]interface{} `yaml:"variables"`
}

func ParseRoot() (*Module, *Override, error) {
	path, err := Path()
	if err != nil {
		return nil, nil, err
	}
	var override Override
	localPath := filepath.Join(filepath.Dir(path), localOverridesFilename)
	if _, err := os.Stat(localPath); err == nil {
		f, err := os.Open(localPath)
		if err != nil {
			return nil, nil, err
		}
		defer f.Close()

		if err := yaml.NewDecoder(f).Decode(&override); err != nil {
			return nil, nil, fmt.Errorf("error in %s: %w", localPath, err)
		}
	} else if !os.IsNotExist(err) {
		return nil, nil, err
	}
	mod, err := Parse(path, BaseVariables())
	if err != nil {
		return nil, nil, err
	}
	return mod, &override, nil
}
