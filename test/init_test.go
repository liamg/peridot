package test

import (
	"testing"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestInitCommand(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Stop() }()

	// run init command
	output, exit, err := c.RunAsUser("peridot", "init")
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

func TestInitCommandCannotOverwriteConfigByDefault(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Stop() }()

	require.NoError(t, c.WriteConfig(``))

	// run init command
	output, exit, err := c.RunAsUser("peridot", "init", "--no-ansi")
	require.NoError(t, err)

	// command should fail when config already exists
	require.Equal(t, 1, exit, output)

}
func TestInitCommandWithForcedOverwrite(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Stop() }()

	// place a config
	require.NoError(t, c.WriteConfig(``))

	// run init command with force
	output, exit, err := c.RunAsUser("peridot", "init", "--force", "--no-ansi")
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

func TestInitCommandWithForcedOverwriteWhenConfigIsInvalid(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Stop() }()

	// place a bad config
	require.NoError(t, c.WriteConfig(`this is invalid`))

	// run init command with force
	output, exit, err := c.RunAsUser("peridot", "init", "--force", "--no-ansi")
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
