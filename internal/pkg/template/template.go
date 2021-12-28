package template

import (
	"io"
	"io/ioutil"
	"text/template"
)

func Apply(r io.Reader, w io.Writer, vars map[string]interface{}) error {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	t, err := template.New("tmp").Parse(string(raw))
	if err != nil {
		return err
	}
	err = t.Execute(w, vars)
	if err != nil {
		return err
	}
	return nil
}
