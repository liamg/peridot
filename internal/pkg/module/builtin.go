package module

import (
	"fmt"
	"os"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/variable"
)

type baseBuiltin struct {
	name      string
	inputs    []config.Variable
	variables variable.Collection
	filesFunc func(variable.Collection) []File
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
	return false
}

func (b *baseBuiltin) RequiresInstall() bool {
	return false
}

func (b *baseBuiltin) Install() error {
	return nil
}

func (b *baseBuiltin) Update() error {
	return nil
}

func (b *baseBuiltin) AfterFileChange() error {
	return nil
}

func (b *baseBuiltin) ApplyVariables(vars variable.Collection) {
	b.variables = applyVariableDefaults(b.inputs, vars)
}
