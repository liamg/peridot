package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

type File struct {
	Name          string            `yaml:"name"`
	Includes      []string          `yaml:"includes"`
	Variables     map[string]string `yaml:"variables"`
	TemplateFiles []TemplateFile    `yaml:"files"`
}

type TemplateFile struct {
	Target string `yaml:"target"`
	Source string `yaml:"source"`
	Raw    bool   `yaml:"raw"`
}

func Parse(r io.Reader) (*File, error) {
	var f File
	if err := yaml.NewDecoder(r).Decode(&f); err != nil {
		return nil, err
	}
	if f.Variables == nil {
		f.Variables = make(map[string]string)
	}
	return &f, nil
}
