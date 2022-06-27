package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveDuplication(t *testing.T) {
	type args struct {
		in []string
	}
	tests := []struct {
		name    string
		args    args
		wantOut []string
	}{
		{
			name:    "test",
			args:    args{in: []string{"1", "2", "4", "4", "1"}},
			wantOut: []string{"1", "2", "4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantOut, RemoveDuplication(tt.args.in), "RemoveDuplication(%v)", tt.args.in)
		})
	}
}
