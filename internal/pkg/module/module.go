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

func (m *module) RequiresUpdate() bool {
	if m.conf.Scripts.UpdateRequired.Command == "" {
		return false
	}
	return run.Run(
		m.conf.Scripts.UpdateRequired.Command,
		m.Path(),
		m.conf.Scripts.UpdateRequired.Sudo,
	) == nil
}

func (m *module) RequiresInstall() bool {
	if m.conf.Scripts.InstallRequired.Command == "" {
		return false
	}
	return run.Run(
		m.conf.Scripts.InstallRequired.Command,
		m.Path(),
		m.conf.Scripts.InstallRequired.Sudo,
	) == nil
}

func (m *module) Install() error {
	return run.Run(
		m.conf.Scripts.Install.Command,
		m.Path(),
		m.conf.Scripts.Install.Sudo,
	)
}

func (m *module) Update() error {
	return run.Run(
		m.conf.Scripts.Update.Command,
		m.Path(),
		m.conf.Scripts.Update.Sudo,
	)
}

func (m *module) AfterFileChange() error {
	if m.conf.Scripts.AfterFileChange.Command == "" {
		return nil
	}
	return run.Run(
		m.conf.Scripts.AfterFileChange.Command,
		m.Path(),
		m.conf.Scripts.AfterFileChange.Sudo,
	)
}

func (m *module) Validate() error {
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
