package printer

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestYamlPrinter(t *testing.T) {
	defer os.RemoveAll("a.yaml")
	printer := NewYamlPrinter()
	assert.NotNil(t, printer)
	err := printer.Print("a.yaml", func() interface{} {
		return struct {
			Name string `yaml:"name"`
		}{
			Name: "123",
		}
	})
	assert.Nil(t, err)
	file, err := ioutil.ReadFile("a.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, file)

}
