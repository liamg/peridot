package template

import (
	"io"
	"io/ioutil"
	"text/template"

	"github.com/liamg/peridot/internal/pkg/variable"
)

func Apply(r io.Reader, w io.Writer, vars variable.Collection) error {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	t, err := template.New("tmp").Parse(string(raw))
	if err != nil {
		return err
	}
	err = t.Execute(w, vars.AsMap())
	if err != nil {
		return err
	}
	return nil
}
