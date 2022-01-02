package test

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystem(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()

	// run system command
	output, exit, err := c.Run("peridot", "system")
	require.NoError(t, err)
	require.Equal(t, 0, exit, output)

	required := map[string]bool{
		"Operating System": false,
		"Architecture":     false,
		"Distribution":     false,
	}

	ansiRegex := regexp.MustCompile(`^.*\[0m`)

	for _, line := range strings.Split(output, "\n") {
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		key = ansiRegex.ReplaceAllString(key, "")
		value := strings.TrimSpace(parts[1])
		value = strings.ReplaceAll(value, "\x1b[0m", "")
		value = strings.ReplaceAll(value, "\x1b[1m", "")
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
