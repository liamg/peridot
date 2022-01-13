package test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemCommand(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Stop() }()

	// run system command
	output, exit, err := c.RunAsUser("peridot", "system", "--no-ansi")
	require.NoError(t, err)
	require.Equal(t, 0, exit, output)

	required := map[string]bool{
		"Operating System": false,
		"Architecture":     false,
		"Distribution":     false,
	}

	for _, line := range strings.Split(output, "\n") {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "Architecture":
			assert.Equal(t, runtime.GOARCH, value)
		case "Operating System":
			assert.Equal(t, runtime.GOOS, value)
		case "Distribution":
			assert.Equal(t, "ubuntu", value)
		default:
			t.Errorf("unexpected key: %q", key)
			continue
		}
		required[key] = true
	}

	for name, ok := range required {
		assert.True(t, ok, fmt.Sprintf("missing property '%s'", name))
	}
}
