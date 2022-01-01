package module

import (
	"testing"

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

	factory.WithRequiresInstallFunc(func() bool {
		requiresInstall = true
		return false
	})
	factory.WithRequiresUpdateFunc(func() bool {
		requiresUpdate = true
		return false
	})
	factory.WithInstallFunc(func() error {
		install = true
		return nil
	})
	factory.WithUpdateFunc(func() error {
		update = true
		return nil
	})
	factory.WithAfterFileChangeFunc(func() error {
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
