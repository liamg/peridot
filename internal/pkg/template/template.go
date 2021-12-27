package template

import (
	"io"
	"io/ioutil"
	"text/template"
)

type Input struct {
	Variables map[string]string
}

func Apply(r io.Reader, w io.Writer, input Input) error {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	t, err := template.New("tmp").Funcs(template.FuncMap{
		"var": func(name string) string {
			if val, ok := input.Variables[name]; ok {
				return val
			}
			return "fuck"
		},
	}).Parse(string(raw))
	if err != nil {
		return err
	}
	err = t.Execute(w, input.Variables)
	if err != nil {
		return err
	}
	return nil
}
