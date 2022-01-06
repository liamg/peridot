package module

import (
	"fmt"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/variable"
)

type Module interface {
	Name() string
	Path() string
	Children() []Module
	Files() []File
	Validate() error
	RequiresUpdate() bool
	RequiresInstall() bool
	Install() error
	Update() error
	AfterFileChange() error
}

type module struct {
	conf      config.Module
	children  []Module
	files     []File
	variables variable.Collection
}

func (m *module) Name() string {
	return m.conf.Name
}

func (m *module) Path() string {
	return m.conf.Dir
}

func (m *module) Children() []Module {
	return m.children
}

func (m *module) Files() []File {
	return m.files
}

func runScript(m Module, s config.Script) error {
	r := NewRunner(m)
	return r.Run(s.Command, s.Sudo)
}

func (m *module) RequiresUpdate() bool {
	if m.conf.Scripts.UpdateRequired.Command == "" {
		return false
	}
	return runScript(
		m,
		m.conf.Scripts.UpdateRequired,
	) == nil
}

func (m *module) RequiresInstall() bool {
	if m.conf.Scripts.InstallRequired.Command == "" {
		return false
	}
	return runScript(
		m,
		m.conf.Scripts.InstallRequired,
	) == nil
}

func (m *module) Install() error {
	return runScript(
		m,
		m.conf.Scripts.Install,
	)
}

func (m *module) Update() error {
	return runScript(
		m,
		m.conf.Scripts.Update,
	)
}

func (m *module) AfterFileChange() error {
	if m.conf.Scripts.AfterFileChange.Command == "" {
		return nil
	}
	return runScript(
		m,
		m.conf.Scripts.AfterFileChange,
	)
}

func (m *module) Validate() error {
	for _, file := range m.conf.Files {
		if file.Target == "" {
			return fmt.Errorf("file has no target")
		}
		if file.Source == "" {
			return fmt.Errorf("file has no source")
		}
	}
	for _, v := range m.conf.Variables {
		if !v.Required {
			continue
		}
		if !m.variables.Has(v.Name) {
			return fmt.Errorf("module '%s' is missing a required variable '%s'", m.Name(), v.Name)
		}
	}
	return nil
}
