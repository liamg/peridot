package variable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCollection(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		want  map[string]interface{}
	}{
		{
			name:  "basic",
			input: map[string]interface{}{},
			want:  map[string]interface{}{},
		},
		{
			name:  "nil",
			input: nil,
			want:  map[string]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coll := NewCollection(tt.input)
			assert.Equal(t, tt.want, coll.AsMap())
		})
	}
}

func Test_collection_Has(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		key   string
		want  bool
	}{
		{
			name:  "does not exist",
			input: map[string]interface{}{},
			key:   "x",
			want:  false,
		},
		{
			name: "exists",
			input: map[string]interface{}{
				"x": "",
			},
			key:  "x",
			want: true,
		},
		{
			name: "different case exists",
			input: map[string]interface{}{
				"x": "",
			},
			key:  "X",
			want: false,
		},
		{
			name: "exists as nil",
			input: map[string]interface{}{
				"x": nil,
			},
			key:  "x",
			want: true,
		},
		{
			name: "exists as map",
			input: map[string]interface{}{
				"x": map[string]interface{}(nil),
			},
			key:  "x",
			want: true,
		},
		{
			name:  "whole map is nil",
			input: nil,
			key:   "x",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coll := NewCollection(tt.input)
			has := coll.Has(tt.key)
			assert.Equal(t, tt.want, has)
		})
	}
}

func Test_collection_Get(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		key   string
		want  Variable
	}{
		{
			name: "exists",
			input: map[string]interface{}{
				"x": 123,
			},
			key:  "x",
			want: New(123),
		},
		{
			name: "does not exist",
			input: map[string]interface{}{
				"x": 123,
			},
			key:  "y",
			want: New(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coll := NewCollection(tt.input)
			actual := coll.Get(tt.key)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func Test_collection_Set(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
		key   string
		val   string
	}{
		{
			name:  "empty set",
			input: nil,
			key:   "x",
			val:   "y",
		},
		{
			name: "overwrite",
			input: map[string]interface{}{
				"hello": "goodbye",
			},
			key: "hello",
			val: "world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coll := NewCollection(tt.input)
			coll.Set(tt.key, tt.val)
			assert.Equal(t, tt.val, coll.AsMap()[tt.key])
		})
	}
}

func Test_collection_MergeIn(t *testing.T) {
	tests := []struct {
		name string
		a    map[string]interface{}
		b    map[string]interface{}
		want map[string]interface{}
	}{
		{
			a: map[string]interface{}{
				"x": 1,
				"y": 2,
			},
			b: map[string]interface{}{
				"x": "one",
				"z": "three",
			},
			want: map[string]interface{}{
				"x": "one",
				"y": 2,
				"z": "three",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewCollection(tt.a)
			b := NewCollection(tt.b)
			a.MergeIn(b)
			assert.Equal(t, tt.want, a.AsMap())
		})
	}
}
