package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModule_Validate(t *testing.T) {
	tests := []struct {
		name    string
		module  Module
		wantErr bool
	}{
		{
			name:    "empty config",
			module:  Module{},
			wantErr: false,
		},
		{
			name: "variable with no name",
			module: Module{
				Variables: []Variable{
					{
						Name: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "variable not required but no default",
			module: Module{
				Variables: []Variable{
					{
						Name:     "x",
						Required: false,
						Default:  nil,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "variable required but has default",
			module: Module{
				Variables: []Variable{
					{
						Name:     "x",
						Required: true,
						Default:  "blah",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "child module without name",
			module: Module{
				Modules: []InnerModule{
					{
						Name: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "child module without source",
			module: Module{
				Modules: []InnerModule{
					{
						Name:   "X",
						Source: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "scripts: has install_required but no install",
			module: Module{
				Scripts: Scripts{
					InstallRequired: "check",
					Install:         "",
				},
			},
			wantErr: true,
		},
		{
			name: "scripts: has no install_required but has install",
			module: Module{
				Scripts: Scripts{
					InstallRequired: "",
					Install:         "check",
				},
			},
			wantErr: true,
		},
		{
			name: "scripts: has update_required but no update",
			module: Module{
				Scripts: Scripts{
					UpdateRequired: "check",
					Update:         "",
				},
			},
			wantErr: true,
		},
		{
			name: "scripts: has no update_required but has update",
			module: Module{
				Scripts: Scripts{
					UpdateRequired: "",
					Update:         "check",
				},
			},
			wantErr: true,
		},
		{
			name: "scripts: has after_file_change but no files",
			module: Module{
				Scripts: Scripts{
					AfterFileChange: "check",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.module.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
