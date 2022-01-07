package module

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"github.com/liamg/peridot/internal/pkg/log"
	"github.com/liamg/tml"
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

	if withContent {
		tml.Printf("<yellow>[Module %s] Changes required for '%s':</yellow>\n", d.module.Name(), d.path)
		edits := myers.ComputeEdits(span.URI(d.path), d.before, d.after)
		for _, line := range strings.Split(fmt.Sprintf("%s", gotextdiff.ToUnified("before", "after", d.before, edits)), "\n") {
			switch {
			case len(line) > 0 && line[0] == '+':
				tml.Printf("<green>%s</green>\n", line)
			case len(line) > 0 && line[0] == '-':
				tml.Printf("<red>%s</red>\n", line)
			default:
				fmt.Println(line)
			}
		}
	} else {
		tml.Printf("<yellow>[Module %s] Changes required for '%s'.</yellow>\n", d.module.Name(), d.path)
	}
}

func (d *fileDiff) Apply() error {

	logger := log.NewLogger(d.module.Name())

	switch d.operation {
	case OpDelete:
		logger.Log("Removing file %s...", d.path)
		return os.Remove(d.path)
	default:
		dir := filepath.Dir(d.path)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
		logger.Log("Writing file %s...", d.path)
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
