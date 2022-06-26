package main

import (
	"github.com/pzs-pzs/cherry-pick/pkg/flow"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"testing"
)

func Test_run(t *testing.T) {
	defer os.RemoveAll("a.yaml")
	err := run("../../testdata.git/", "a.yaml")
	assert.Nil(t, err)
	file, err := ioutil.ReadFile("a.yaml")
	assert.Nil(t, err)
	var rst []*flow.PrintData
	err = yaml.Unmarshal(file, &rst)
	assert.Nil(t, err)
	for _, r := range rst {
		assert.Equal(t, "8fe8f231cf539e3346a4fd31d9c275bf168f6cc8", r.OriginalCommit)
		assert.Equal(t, []string{"c8c9bb7032de90f8d3e9107c6c2aa569855f4796"}, r.CherryPickCommits)
	}
}
