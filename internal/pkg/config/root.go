package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Dir   string `yaml:"-"`
	Path  string `yaml:"-"`
	Debug bool   `yaml:"debug"`
}

func ParseRoot() (*Config, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}
	c.Dir = filepath.Dir(path)
	c.Path = path
	return &c, nil
}
