package module

import (
	"fmt"
	"os"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type FileOperation uint8

const (
	OpUnknown FileOperation = iota
	OpCreate
	OpUpdate
	OpDelete
)

type fileDiff struct {
	module    Module
	path      string
	before    string
	after     string
	operation FileOperation
}

func (d *fileDiff) Module() Module {
	return d.module
}

func (d *fileDiff) Path() string {
	return d.path
}

func (d *fileDiff) Before() string {
	return d.before
}

func (d *fileDiff) After() string {
	return d.after
}

func (d *fileDiff) Operation() FileOperation {
	return d.operation
}

func (d *fileDiff) Print(withContent bool) {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(d.before, d.after, false)

	if withContent {
		fmt.Printf("File '%s' has pending changes:\n", d.path)
		fmt.Println(dmp.DiffPrettyText(diffs))
	} else {
		fmt.Printf("File '%s' has pending changes.\n", d.path)
	}
}

func (d *fileDiff) Apply() error {
	switch d.operation {
	case OpDelete:
		return os.Remove(d.path)
	default:
		f, err := os.Create(d.path)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()
		if _, err := f.WriteString(d.after); err != nil {
			return err
		}
	}

	return nil
}
