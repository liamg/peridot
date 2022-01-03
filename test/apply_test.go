package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyCommandWithEmptyConfig(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()

	require.NoError(t, c.WriteConfig(``))

	output, exit, err := c.Run("peridot", "apply", "--no-ansi")
	require.NoError(t, err)
	assert.Equal(t, 0, exit, output)
	assert.Contains(t, output, "no changes")
}

func TestApplyCommandWithInvalidConfig(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()

	require.NoError(t, c.WriteConfig(`this is invalid`))

	output, exit, err := c.Run("peridot", "apply", "--no-ansi")
	require.NoError(t, err)
	assert.Equal(t, 1, exit, output)
}

func TestApplyCommandWithSingleFileChanges(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()

	require.NoError(t, c.WriteHomeFile(".config/peridot/lol.tmpl", `hello world`))

	require.NoError(t, c.WriteConfig(`
files:
  - target: "{{ .user_home_dir }}/hello.txt"
    template: ./lol.tmpl
`))

	output, exit, err := c.Run("peridot", "apply", "--no-ansi")
	require.NoError(t, err)
	assert.Equal(t, 0, exit, output)

	actual, err := c.ReadHomeFile("hello.txt")
	require.NoError(t, err)
	assert.Equal(t, "hello world", actual)
}

func TestApplyCommandWhenOnlyFileAlreadyMatches(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()

	require.NoError(t, c.WriteHomeFile(".config/peridot/lol.tmpl", `hello world`))

	require.NoError(t, c.WriteHomeFile("hello.txt", `hello world`))

	require.NoError(t, c.WriteConfig(`
files:
  - target: "{{ .user_home_dir }}/hello.txt"
    template: ./lol.tmpl
`))

	output, exit, err := c.Run("peridot", "apply", "--no-ansi")
	require.NoError(t, err)
	assert.Equal(t, 0, exit, output)
	assert.Contains(t, output, "no changes")
}
