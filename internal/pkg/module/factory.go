package module

import (
	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/variable"
)

type factory struct {
	base baseBuiltin
}

func NewFactory(name string) *factory {
	return &factory{
		base: baseBuiltin{
			name: name,
		},
	}
}

func (f *factory) Build() BuiltIn {
	return &f.base
}

func (f *factory) WithInputs(inputs []config.Variable) *factory {
	f.base.inputs = inputs
	return f
}

func (f *factory) WithFilesFunc(fnc func(vars variable.Collection) []File) *factory {
	f.base.filesFunc = fnc
	return f
}

func (f *factory) WithRequiresUpdateFunc(fnc func(*Runner, variable.Collection) (bool, error)) *factory {
	f.base.requiresUpdateFunc = fnc
	return f
}

func (f *factory) WithRequiresInstallFunc(fnc func(*Runner, variable.Collection) (bool, error)) *factory {
	f.base.requiresInstallFunc = fnc
	return f
}

func (f *factory) WithUpdateFunc(fnc func(*Runner, variable.Collection) error) *factory {
	f.base.updateFunc = fnc
	return f
}

func (f *factory) WithInstallFunc(fnc func(*Runner, variable.Collection) error) *factory {
	f.base.installFunc = fnc
	return f
}

func (f *factory) WithAfterFileChangeFunc(fnc func(*Runner, variable.Collection) error) *factory {
	f.base.afterFileChangeFunc = fnc
	return f
}
