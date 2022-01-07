package module

import (
	"fmt"
	"os"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/log"
	"github.com/liamg/peridot/internal/pkg/variable"
)

type baseBuiltin struct {
	name                string
	inputs              []config.Variable
	variables           variable.Collection
	filesFunc           func(variable.Collection) []File
	requiresInstallFunc func(*Runner, variable.Collection) (bool, error)
	requiresUpdateFunc  func(*Runner, variable.Collection) (bool, error)
	installFunc         func(*Runner, variable.Collection) error
	updateFunc          func(*Runner, variable.Collection) error
	afterFileChangeFunc func(*Runner, variable.Collection) error
}

func (b *baseBuiltin) Name() string {
	return b.name
}

func (b *baseBuiltin) Clone(name string) BuiltIn {
	if b == nil {
		return nil
	}
	c := *b
	c.name = name
	return &c
}

func (b *baseBuiltin) Path() string {
	return os.TempDir()
}

func (b *baseBuiltin) Children() []Module {
	return nil
}

func (b *baseBuiltin) Files() []File {
	if b.filesFunc == nil {
		return nil
	}
	return b.filesFunc(b.variables)
}

func (b *baseBuiltin) Validate() error {
	return validateVariables(b.inputs, b.variables)
}

func (b *baseBuiltin) RequiresUpdate() bool {
	if b.requiresUpdateFunc == nil {
		return false
	}

	required, err := b.requiresUpdateFunc(NewRunner(b, "should_update"), b.variables)
	if err != nil {
		log.Debug("[%s] Non-zero exit checking if update was required: %s", b.Name(), err)
	}
	return required
}

func (b *baseBuiltin) RequiresInstall() bool {
	if b.requiresInstallFunc == nil {
		return false
	}
	required, err := b.requiresInstallFunc(NewRunner(b, "should_install"), b.variables)
	if err != nil {
		log.Debug("[%s] Non-zero exit checking if install was required: %s", b.Name(), err)
	}
	return required
}

func (b *baseBuiltin) Install() error {
	if b.installFunc == nil {
		return fmt.Errorf("install handler not implemented")
	}
	return b.installFunc(NewRunner(b, "install"), b.variables)
}

func (b *baseBuiltin) Update() error {
	if b.updateFunc == nil {
		return fmt.Errorf("update handler not implemented")
	}
	return b.updateFunc(NewRunner(b, "update"), b.variables)
}

func (b *baseBuiltin) AfterFileChange() error {
	if b.afterFileChangeFunc == nil {
		return nil
	}
	return b.afterFileChangeFunc(NewRunner(b, "after_file_change"), b.variables)
}

func (b *baseBuiltin) ApplyVariables(vars variable.Collection) {
	b.variables = applyVariableDefaults(b.inputs, vars)
}
