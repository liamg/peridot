package variable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList_All(t *testing.T) {
	tests := []struct {
		name string
		l    List
		want []interface{}
	}{
		{
			name: "empty",
			l:    make(List, 0),
			want: []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.l.All()
			assert.Equal(t, len(tt.want), len(actual))
			for i, a := range actual {
				assert.Equal(t, tt.want[i], a.Interface())
			}
		})
	}
}

func TestList_Len(t *testing.T) {
	tests := []struct {
		name string
		l    List
		want int
	}{
		{
			name: "empty list",
			l:    make(List, 0),
			want: 0,
		},
		{
			name: "non-empty list",
			l:    make(List, 1337),
			want: 1337,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Len(); got != tt.want {
				t.Errorf("List.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
