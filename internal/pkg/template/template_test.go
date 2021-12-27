package template

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariables(t *testing.T) {

	raw := `hello {{ var "test" }}`
	input := Input{
		Variables: map[string]string{
			"test": "world",
		},
	}
	var output []byte
	buffer := bytes.NewBuffer(output)
	err := Apply(strings.NewReader(raw), buffer, input)
	require.NoError(t, err)

	assert.Equal(t, `hello world`, buffer.String())

}
