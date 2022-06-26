package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_run(t *testing.T) {
	defer os.RemoveAll("./a.yaml")
	err := run("../../../pick", "./a.yaml")
	assert.Nil(t, err)
}
