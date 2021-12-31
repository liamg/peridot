package module

import (
	"fmt"
	"os"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/variable"
)

type baseBuiltin struct {
	name                string
	inputs              []config.Variable
	variables           variable.Collection
	filesFunc           func(variable.Collection) []File
	requiresInstallFunc func() bool
	requiresUpdateFunc  func() bool
	installFunc         func() error
	updateFunc          func() error
	afterFileChangeFunc func() error
}

func (b *baseBuiltin) Name() string {
	return fmt.Sprintf("builtin:%s", b.name)
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
	return b.requiresUpdateFunc()
}

func (b *baseBuiltin) RequiresInstall() bool {
	if b.requiresInstallFunc == nil {
		return false
	}
	return b.requiresInstallFunc()
}

func (b *baseBuiltin) Install() error {
	if b.installFunc == nil {
		return fmt.Errorf("install handler not implemented")
	}
	return b.installFunc()
}

func (b *baseBuiltin) Update() error {
	if b.updateFunc == nil {
		return fmt.Errorf("update handler not implemented")
	}
	return b.updateFunc()
}

func (b *baseBuiltin) AfterFileChange() error {
	if b.afterFileChangeFunc == nil {
		return nil
	}
	return b.afterFileChangeFunc()
}

func (b *baseBuiltin) ApplyVariables(vars variable.Collection) {
	b.variables = applyVariableDefaults(b.inputs, vars)
}
