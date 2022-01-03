package module

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/template"
	"github.com/liamg/peridot/internal/pkg/variable"
)

type File interface {
	Target() string
	RenderTemplate() (string, error)
}

func NewMemoryFile(target string, template string, enableTemplating bool, vars variable.Collection) File {
	return &memoryFile{
		target:           target,
		template:         template,
		variables:        vars,
		enableTemplating: enableTemplating,
	}
}

func (m *memoryFile) Target() string {
	return m.target
}

func (m *memoryFile) RenderTemplate() (string, error) {
	if !m.enableTemplating {
		return m.template, nil
	}
	buffer := bytes.NewBufferString("")
	if err := template.Apply(strings.NewReader(m.template), buffer, m.variables); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

type memoryFile struct {
	target           string
	template         string
	variables        variable.Collection
	enableTemplating bool
}

type localFile struct {
	target     string
	sourcePath string
	variables  variable.Collection
	templating bool
}

func loadFile(modConf config.Module, fileConf config.File, combined variable.Collection) (File, error) {
	templatePath := filepath.Join(modConf.Dir, fileConf.Source)
	return &localFile{
		target:     fileConf.Target,
		sourcePath: templatePath,
		variables:  combined,
		templating: !fileConf.DisableTemplating,
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
	if !l.templating {
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	buffer := bytes.NewBufferString("")
	if err := template.Apply(f, buffer, l.variables); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
