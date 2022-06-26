package printer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYamlPrinter(t *testing.T) {
	printer := NewYamlPrinter()
	assert.NotNil(t, printer)
	err := printer.Print("./a.yaml", func() interface{} {
		return struct {
			Name string `yaml:"name"`
		}{
			Name: "123",
		}
	})
	assert.Nil(t, err)
}
