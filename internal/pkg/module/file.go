package module

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/template"
)

type File interface {
	Target() string
	RenderTemplate() (string, error)
}

func NewMemoryFile(target string, template string, vars map[string]interface{}) File {
	return &memoryFile{
		target:    target,
		template:  template,
		variables: vars,
	}
}

func (m *memoryFile) Target() string {
	return m.target
}

func (m *memoryFile) RenderTemplate() (string, error) {
	buffer := bytes.NewBufferString("")
	if err := template.Apply(strings.NewReader(m.template), buffer, m.variables); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

type memoryFile struct {
	target    string
	template  string
	variables map[string]interface{}
}

type localFile struct {
	target     string
	sourcePath string
	variables  map[string]interface{}
}

func loadFile(modConf config.Module, fileConf config.File, combined map[string]interface{}) (File, error) {
	templatePath := filepath.Join(modConf.Dir, fileConf.Template)
	return &localFile{
		target:     fileConf.Target,
		sourcePath: templatePath,
		variables:  combined,
	}, nil
}

func (l *localFile) Target() string {
	return l.target
}

func (l *localFile) RenderTemplate() (string, error) {
	f, err := os.Open(l.sourcePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	buffer := bytes.NewBufferString("")
	if err := template.Apply(f, buffer, l.variables); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
