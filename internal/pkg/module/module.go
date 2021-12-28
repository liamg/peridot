package module

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/liamg/peridot/internal/pkg/config"
)

type Module interface {
	Name() string
	Path() string
	Children() []Module
	Files() []File
	Diff() ([]ModuleDiff, []FileDiff, error)
	Validate() error
	RequiresUpdate() bool
	RequiresInstall() bool
	Install() error
	Update() error
}

type module struct {
	conf      config.Module
	children  []Module
	files     []File
	variables map[string]interface{}
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

func (m *module) Diff() ([]ModuleDiff, []FileDiff, error) {
	var fileDiffs []FileDiff
	var moduleDiffs []ModuleDiff

	// run scripts.update_required and scripts.install_required to see if update is needed
	if m.RequiresInstall() {
		moduleDiffs = append(moduleDiffs, &moduleDiff{
			module: m,
			before: StateUninstalled,
			after:  StateInstalled,
		})
	} else if m.RequiresUpdate() {
		moduleDiffs = append(moduleDiffs, &moduleDiff{
			module: m,
			before: StateInstalled,
			after:  StateUpdated,
		})
	}

	for _, mod := range m.Children() {
		m, f, err := mod.Diff()
		if err != nil {
			return nil, nil, err
		}
		moduleDiffs = append(moduleDiffs, m...)
		fileDiffs = append(fileDiffs, f...)
	}

	for _, file := range m.Files() {
		if err := func() error {
			diff := fileDiff{
				module:    m,
				path:      file.Target(),
				operation: OpCreate,
			}
			targetFile, err := os.Open(file.Target())
			if err == nil {
				content, err := ioutil.ReadAll(targetFile)
				if err != nil {
					return err
				}
				_ = targetFile.Close()
				diff.before = string(content)
				diff.operation = OpUpdate
			}

			after, err := file.RenderTemplate()
			if err != nil {
				return err
			}
			diff.after = after
			fileDiffs = append(fileDiffs, &diff)
			return nil
		}(); err != nil {
			return nil, nil, err
		}
	}
	return moduleDiffs, fileDiffs, nil
}

func (m *module) Validate() error {
	for _, v := range m.conf.Variables {
		if !v.Required {
			continue
		}
		if _, ok := m.variables[v.Name]; !ok {
			return fmt.Errorf("module '%s' is missing a required variable '%s'", m.Name(), v.Name)
		}
	}
	return nil
}
