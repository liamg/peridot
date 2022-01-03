package module

import (
	"testing"

	"github.com/liamg/peridot/internal/pkg/variable"
	"github.com/stretchr/testify/assert"
)

func TestFactorySetsModuleName(t *testing.T) {
	factory := NewFactory("example")
	built := factory.Build()
	assert.Equal(t, "example", built.Name())
}

func TestFactorySetsHandlers(t *testing.T) {

	factory := NewFactory("example")

	var requiresInstall, requiresUpdate, install, update, after bool

	factory.WithRequiresInstallFunc(func(variable.Collection) bool {
		requiresInstall = true
		return false
	})
	factory.WithRequiresUpdateFunc(func(variable.Collection) bool {
		requiresUpdate = true
		return false
	})
	factory.WithInstallFunc(func(variable.Collection) error {
		install = true
		return nil
	})
	factory.WithUpdateFunc(func(variable.Collection) error {
		update = true
		return nil
	})
	factory.WithAfterFileChangeFunc(func(variable.Collection) error {
		after = true
		return nil
	})

	built := factory.Build()

	built.RequiresInstall()
	built.RequiresUpdate()
	built.Install()
	built.Update()
	built.AfterFileChange()

	assert.True(t, requiresInstall, "RequiresInstall() was not configured")
	assert.True(t, requiresUpdate, "RequiresUpdate() was not configured")
	assert.True(t, install, "Install() was not configured")
	assert.True(t, update, "Update() was not configured")
	assert.True(t, after, "AfterFileChange() was not configured")
}
