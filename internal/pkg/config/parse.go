package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/liamg/peridot/internal/pkg/template"
	"github.com/liamg/peridot/internal/pkg/variable"
	"gopkg.in/yaml.v3"
)

func ParseVariables(path string) ([]Variable, error) {
	var sniff variableSniff
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	if err := yaml.NewDecoder(f).Decode(&sniff); err != nil {
		return nil, fmt.Errorf("error in %s: %w", path, err)
	}
	return sniff.Variables, nil
}

func Parse(path string, variables variable.Collection) (*Module, error) {
	var m Module
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var output []byte
	buffer := bytes.NewBuffer(output)
	if err := template.Apply(f, buffer, variables); err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(buffer.Bytes(), &m); err != nil {
		return nil, fmt.Errorf("error in %s: %w", path, err)
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	m.Path = path
	m.Dir = filepath.Dir(path)
	return &m, nil
}
