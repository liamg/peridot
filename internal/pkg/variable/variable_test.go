package variable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "bool",
			input: true,
		},
		{
			name:  "int",
			input: 123,
		},
		{
			name:  "string",
			input: "hello",
		},
		{
			name:  "nil",
			input: nil,
		},
		{
			name:  "float",
			input: 12.3,
		},
		{
			name:  "[]string",
			input: []string{"a", "b", "c"},
		},
		{
			name:  "[]int",
			input: []int{123, 456},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(tt.input)
			assert.Equal(t, tt.input, v.Interface())
		})
	}
}

func Test_variable_AsString(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{
			name:  "string",
			input: "hello world",
			want:  "hello world",
		},
		{
			name:  "bool",
			input: true,
			want:  "true",
		},
		{
			name:  "int",
			input: 123,
			want:  "123",
		},
		{
			name:  "float",
			input: 1.234567,
			want:  "1.234567",
		},
		{
			name:  "nil",
			input: nil,
			want:  "<nil>",
		},
		{
			name:  "[]int",
			input: []int{1, 2, 3},
			want:  "[1, 2, 3]",
		},
		{
			name:  "[]float64",
			input: []float64{1.1, 2.2, 3.3},
			want:  "[1.100000, 2.200000, 3.300000]",
		},
		{
			name:  "[]string",
			input: []string{"a", "b", "c"},
			want:  `["a", "b", "c"]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(tt.input)
			assert.Equal(t, tt.want, v.AsString())
		})
	}
}

func Test_variable_AsInteger(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  int
	}{
		{
			name:  "int",
			input: 123,
			want:  123,
		},
		{
			name:  "bool false",
			input: false,
			want:  0,
		},
		{
			name:  "bool true",
			input: true,
			want:  1,
		},
		{
			name:  "string",
			input: "123",
			want:  123,
		},
		{
			name:  "float",
			input: 1.23,
			want:  1,
		},
		{
			name:  "nil",
			input: nil,
			want:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(tt.input)
			assert.Equal(t, tt.want, v.AsInteger())
		})
	}
}

func Test_variable_AsBool(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  bool
	}{
		{
			name:  "bool true",
			input: true,
			want:  true,
		},
		{
			name:  "bool false",
			input: false,
			want:  false,
		},
		{
			name:  "int == 0",
			input: 0,
			want:  false,
		},
		{
			name:  "int > 0",
			input: 123,
			want:  true,
		},
		{
			name:  "string empty",
			input: "",
			want:  false,
		},
		{
			name:  "string non-empty",
			input: "lol",
			want:  true,
		},
		{
			name:  "nil",
			input: nil,
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(tt.input)
			assert.Equal(t, tt.want, v.AsBool())
		})
	}
}

func Test_variable_AsFloat(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  float64
	}{
		{
			name:  "float",
			input: 1.23,
			want:  1.23,
		},
		{
			name:  "bool false",
			input: false,
			want:  0.0,
		},
		{
			name:  "bool true",
			input: true,
			want:  1.0,
		},
		{
			name:  "int",
			input: 123,
			want:  123.0,
		},
		{
			name:  "string",
			input: "1.230",
			want:  1.23,
		},
		{
			name:  "nil",
			input: nil,
			want:  0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(tt.input)
			assert.Equal(t, tt.want, v.AsFloat64())
		})
	}
}

func Test_variable_AsList(t *testing.T) {
	tests := []struct {
		name       string
		input      interface{}
		wantLen    int
		wantValues []interface{}
	}{
		{
			name: "[]string",
			input: []string{
				"a",
				"b",
				"c",
			},
			wantLen: 3,
			wantValues: []interface{}{
				"a",
				"b",
				"c",
			},
		},
		{
			name: "[]int",
			input: []int{
				5,
				6,
				7,
			},
			wantLen: 3,
			wantValues: []interface{}{
				5,
				6,
				7,
			},
		},
		{
			name: "[]float64",
			input: []float64{
				5.1,
				6.1,
				7.1,
			},
			wantLen: 3,
			wantValues: []interface{}{
				5.1,
				6.1,
				7.1,
			},
		},
		{
			name: "[]bool",
			input: []bool{
				true,
				false,
				true,
			},
			wantLen: 3,
			wantValues: []interface{}{
				true,
				false,
				true,
			},
		},
		{
			name: "[]interface{}",
			input: []interface{}{
				"a",
				123,
				"c",
			},
			wantLen: 3,
			wantValues: []interface{}{
				"a",
				123,
				"c",
			},
		},
		{
			name:       "nil",
			input:      nil,
			wantLen:    0,
			wantValues: []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(tt.input)
			assert.Equal(t, tt.wantLen, v.AsList().Len())
			for i, val := range v.AsList().All() {
				assert.Equal(t, tt.wantValues[i], val.Interface())
			}
		})
	}
}

func Test_variable_AsCollection(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "msi",
			input: map[string]interface{}{"hello": "world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New(tt.input)
			assert.Equal(t, tt.input, v.AsCollection().AsMap())
		})
	}
}
