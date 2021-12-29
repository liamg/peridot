package module

import (
	"fmt"
	"os/exec"

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

func (m *module) RequiresUpdate() bool {
	if m.conf.Scripts.UpdateRequired == "" {
		return false
	}
	cmd := exec.Command("sh", "-c", m.conf.Scripts.UpdateRequired)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (m *module) RequiresInstall() bool {
	if m.conf.Scripts.InstallRequired == "" {
		return false
	}
	cmd := exec.Command("sh", "-c", m.conf.Scripts.InstallRequired)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func (m *module) Install() error {
	return exec.Command("sh", "-c", m.conf.Scripts.Install).Run()
}

func (m *module) Update() error {
	return exec.Command("sh", "-c", m.conf.Scripts.Update).Run()
}

func (m *module) AfterFileChange() error {
	if m.conf.Scripts.AfterFileChange == "" {
		return nil
	}
	return exec.Command("sh", "-c", m.conf.Scripts.AfterFileChange).Run()
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
