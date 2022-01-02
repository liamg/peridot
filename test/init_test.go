package test

import (
	"testing"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestInit(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()

	// run init command
	output, exit, err := c.Run("peridot", "init")
	require.NoError(t, err)
	require.Equal(t, 0, exit, output)

	// read in the config file that the init created
	content, err := c.ReadHomeFile(".config/peridot/config.yml")
	require.NoError(t, err)

	// make sure the yaml is valid
	var out config.Module
	require.NoError(t, yaml.Unmarshal([]byte(content), &out))
	assert.NoError(t, out.Validate())

	// make sure we don't have any crap in here by default
	assert.Len(t, out.Files, 0)
	assert.Len(t, out.Modules, 0)
	assert.Len(t, out.Variables, 0)
}

func TestInitCannotOverwriteConfigByDefault(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()

	_, exit, err := c.Run("mkdir", "-p", ".config/peridot")
	require.NoError(t, err)
	require.Equal(t, 0, exit)
	_, exit, err = c.Run("touch", ".config/peridot/config.yml")
	require.NoError(t, err)
	require.Equal(t, 0, exit)

	// run init command
	output, exit, err := c.Run("peridot", "init")
	require.NoError(t, err)

	// command should fail when config already exists
	require.Equal(t, 1, exit, output)

}
func TestInitWithForcedOverwrite(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()

	// place a config
	_, exit, err := c.Run("mkdir", "-p", ".config/peridot")
	require.NoError(t, err)
	require.Equal(t, 0, exit)
	_, exit, err = c.Run("touch", ".config/peridot/config.yml")
	require.NoError(t, err)
	require.Equal(t, 0, exit)

	// run init command with force
	output, exit, err := c.Run("peridot", "init", "--force")
	require.NoError(t, err)
	require.Equal(t, 0, exit, output)

	// read in the config file that the init created
	content, err := c.ReadHomeFile(".config/peridot/config.yml")
	require.NoError(t, err)
	require.Greater(t, len(content), 0)

	// make sure the yaml is valid
	var out config.Module
	require.NoError(t, yaml.Unmarshal([]byte(content), &out))
	assert.NoError(t, out.Validate())

}
