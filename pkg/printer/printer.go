package printer

import (
	"github.com/pkg/errors"
	"github.com/pzs-pzs/cherry-pick/pkg/util"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Printer interface {
	Print(path string, l LoadData) error
}

func NewYamlPrinter() Printer {
	return &yamlPrinter{}
}

type LoadData func() interface{}

type yamlPrinter struct {
}

func (y *yamlPrinter) Print(path string, l LoadData) error {
	if path == "" {
		return errors.New("path is empty,plz check")
	}

	exists, err := util.PathExists(path)
	if err != nil {
		return err
	}
	if exists {
		return errors.Errorf("[%s] exist", path)
	}

	data := l()
	if data == nil {
		return errors.New("print data is nil")
	}
	out, err := yaml.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	err = ioutil.WriteFile(path, out, 0777)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
