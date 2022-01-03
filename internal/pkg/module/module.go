package module

import (
	"fmt"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/run"
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

func runScript(s config.Script, path string) error {
	return run.Run(s.Command, path, s.Sudo, s.Interactive)
}

func (m *module) RequiresUpdate() bool {
	if m.conf.Scripts.UpdateRequired.Command == "" {
		return false
	}
	return runScript(
		m.conf.Scripts.UpdateRequired,
		m.Path(),
	) == nil
}

func (m *module) RequiresInstall() bool {
	if m.conf.Scripts.InstallRequired.Command == "" {
		return false
	}
	return runScript(
		m.conf.Scripts.InstallRequired,
		m.Path(),
	) == nil
}

func (m *module) Install() error {
	return runScript(
		m.conf.Scripts.Install,
		m.Path(),
	)
}

func (m *module) Update() error {
	return runScript(
		m.conf.Scripts.Update,
		m.Path(),
	)
}

func (m *module) AfterFileChange() error {
	if m.conf.Scripts.AfterFileChange.Command == "" {
		return nil
	}
	return runScript(
		m.conf.Scripts.AfterFileChange,
		m.Path(),
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
