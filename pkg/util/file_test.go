package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PathNotExists(t *testing.T) {
	exists, err := PathExists("a.yaml")
	assert.Nil(t, err)
	assert.Equal(t, false, exists)
}

func Test_PathExists(t *testing.T) {
	exists, err := PathExists("file.go")
	assert.Nil(t, err)
	assert.Equal(t, true, exists)
}
