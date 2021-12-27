package module

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/liamg/peridot/internal/pkg/template"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type FileDiff struct {
	Module *Module
	Path   string
	Before string
	After  string
}

func (m *Module) Diff() ([]FileDiff, error) {
	var diffs []FileDiff
	for _, file := range m.Files {
		if err := func() error {
			diff := FileDiff{
				Module: m,
				Path:   file.Target,
			}
			targetFile, err := os.Open(file.Target)
			if err == nil {
				content, err := ioutil.ReadAll(targetFile)
				if err != nil {
					return err
				}
				_ = targetFile.Close()
				diff.Before = string(content)
			}

			srcFile, err := os.Open(file.Source)
			if err != nil {
				return err
			}
			defer srcFile.Close()

			after := bytes.NewBufferString("")
			if err := template.Apply(srcFile, after, template.Input{}); err != nil {
				return err
			}
			diff.After = after.String()
			diffs = append(diffs, diff)
			return nil
		}(); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (d *FileDiff) Print() {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(d.Before, d.After, false)

	fmt.Println(dmp.DiffPrettyText(diffs))
}
