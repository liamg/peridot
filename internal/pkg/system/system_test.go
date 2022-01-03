package system

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	info := Info()
	assert.Equal(t, runtime.GOOS, info.OperatingSystem)
	assert.Equal(t, runtime.GOARCH, info.Architecture)
}
