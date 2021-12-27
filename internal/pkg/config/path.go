package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	dir      = "peridot"
	filename = "config.yml"
)

func Path() (string, error) {
	root, err := configRoot()
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Join(root, dir, filename))
}

func configRoot() (string, error) {
	if root := os.Getenv("XDG_CONFIG_HOME"); root != "" {
		return root, nil
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config"), nil
	}
	return "", fmt.Errorf("cloud not find config directory")
}
