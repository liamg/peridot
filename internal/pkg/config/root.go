package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/liamg/peridot/internal/pkg/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Dir   string `yaml:"-"`
	Path  string `yaml:"-"`
	Debug bool   `yaml:"debug"`
}

type Override struct {
	Variables map[string]interface{} `yaml:"variables"`
}

func ParseRoot() (*Module, *Override, error) {
	logger := log.NewLogger("root")
	path, err := Path()
	if err != nil {
		return nil, nil, err
	}
	logger.Log("Using config path: %s", path)
	var override Override
	localPath := filepath.Join(filepath.Dir(path), localOverridesFilename)
	if _, err := os.Stat(localPath); err == nil {
		logger.Log("Found local overrides: %s", localPath)
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
	base := BaseVariables()

	for key, val := range base.AsMap() {
		logger.Log("Base variable %s => %s", key, val)
	}

	mod, err := Parse(path, base)
	if err != nil {
		return nil, nil, err
	}

	return mod, &override, nil
}
