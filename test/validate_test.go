package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateCommand(t *testing.T) {

	c, err := startContainer("ubuntu")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Stop() }()

	tests := []struct {
		name    string
		config  string
		wantErr bool
	}{
		{
			name:    "empty config is valid",
			config:  ``,
			wantErr: false,
		},
		{
			name:    "invalid config",
			config:  `blah`,
			wantErr: true,
		},
		{
			name: "module without name",
			config: `
modules:
  - source: builtin:git
    
`,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			require.NoError(t, c.WriteConfig(test.config))

			output, exit, err := c.RunAsUser("peridot", "validate", "--no-ansi")
			require.NoError(t, err)
			if test.wantErr {
				require.NotEqual(t, 0, exit, output)
			} else {
				require.Equal(t, 0, exit, output)
			}
		})
	}
}
