package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiffCommandWithEmptyConfig(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Stop() }()

	require.NoError(t, c.WriteConfig(``))

	output, exit, err := c.RunAsUser("peridot", "diff", "--no-ansi")
	require.NoError(t, err)
	assert.Equal(t, 0, exit, output)
	assert.Contains(t, output, "no changes")
}

func TestDiffCommandWithInvalidConfig(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Stop() }()

	require.NoError(t, c.WriteConfig(`this is invalid`))

	output, exit, err := c.RunAsUser("peridot", "diff", "--no-ansi")
	require.NoError(t, err)
	assert.Equal(t, 1, exit, output)
}

func TestDiffCommandWithSingleFileChanges(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Stop() }()

	require.NoError(t, c.WriteHomeFile(".config/peridot/lol.tmpl", `hello world`))

	require.NoError(t, c.WriteConfig(`
files:
  - target: /tmp/lol
    source: ./lol.tmpl
`))

	output, exit, err := c.RunAsUser("peridot", "diff", "--no-ansi")
	require.NoError(t, err)
	assert.Equal(t, 0, exit, output)
	assert.Contains(t, output, "1 module has pending changes")
	assert.Contains(t, output, "'/tmp/lol'")
}

func TestDiffCommandWhenOnlyFileAlreadyMatches(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Stop() }()

	require.NoError(t, c.WriteHomeFile(".config/peridot/lol.tmpl", `hello world`))

	require.NoError(t, c.WriteHomeFile("hello.txt", `hello world`))

	require.NoError(t, c.WriteConfig(`
files:
  - target: "{{ .user_home_dir }}/hello.txt"
    source: ./lol.tmpl
`))

	output, exit, err := c.RunAsUser("peridot", "diff", "--no-ansi")
	require.NoError(t, err)
	assert.Equal(t, 0, exit, output)
	assert.Contains(t, output, "no changes")
}
