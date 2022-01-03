package test

import (
	"fmt"
	"strings"
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

	output, exit, err := c.RunAsUser("peridot", "apply", "--no-ansi")
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

	output, exit, err := c.RunAsUser("peridot", "apply", "--no-ansi")
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

	output, exit, err := c.RunAsUser("peridot", "apply", "--no-ansi")
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

	output, exit, err := c.RunAsUser("peridot", "apply", "--no-ansi")
	require.NoError(t, err)
	assert.Equal(t, 0, exit, output)
	assert.Contains(t, output, "no changes")
}

func TestApplyWithSudoRequired(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()

	// intall sudo and allow default user to sudo without password
	_, exit, err := c.RunAsRoot("sh", "-c", fmt.Sprintf(`apt update && apt install -y sudo && echo '%s ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers`, defaultUser))
	require.NoError(t, err)
	require.Zero(t, exit)

	require.NoError(t, c.WriteHomeFile(
		".config/peridot/sudo/config.yml",
		`scripts:
  should_install:
    command: true
  install:
    command: "echo ok > /installed && chmod 777 /installed"
    sudo: true
`))

	require.NoError(t, c.WriteConfig(
		`modules:
  - name: sudo
    source: ./sudo
`))

	output, exit, err := c.RunAsUser("peridot", "apply", "--no-ansi")
	require.NoError(t, err)
	require.Equal(t, 0, exit, output)

	output, exit, err = c.RunAsUser("cat", "/installed")
	require.NoError(t, err)
	require.Zero(t, exit, output)
	assert.Equal(t, "ok", strings.TrimSpace(output))

	output, exit, err = c.RunAsUser("stat", "-c", "%U", "/installed")
	require.NoError(t, err)
	require.Zero(t, exit, output)
	assert.Equal(t, "root", strings.TrimSpace(output))

}
