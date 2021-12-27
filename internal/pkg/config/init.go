package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var DefaultConfig = Config{
	Debug: false,
}

func Init() (string, error) {
	path, err := Path()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := f.Write([]byte("# See https://github.com/liamg/peridot for help configuring peridot\n")); err != nil {
		return "", err
	}
	if err := yaml.NewEncoder(f).Encode(DefaultConfig); err != nil {
		return "", err
	}
	return path, nil
}
