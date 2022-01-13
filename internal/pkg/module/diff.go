package module

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/liamg/peridot/internal/pkg/log"
)

type State uint8

const (
	StateUnknown State = iota
	StateUninstalled
	StateInstalled
	StateUpdated
)

type ModuleDiff interface {
	Module() Module
	Before() State
	After() State
	Print(withContent bool)
	Apply() error
	Files() []FileDiff
}

type FileDiff interface {
	Module() Module
	Path() string
	Operation() FileOperation
	Before() string
	After() string
	Print(withContent bool)
	Apply() error
}

func Diff(m Module) ([]ModuleDiff, error) {
	var fileDiffs []FileDiff
	var moduleDiffs []ModuleDiff

	logger := log.NewLogger(m.Name())
	logger.Log("Checking module for required changes...")

	for _, file := range m.Files() {
		logger.Log("Comparing %s...", file.Target())
		if err := func() error {
			diff := fileDiff{
				module:    m,
				path:      file.Target(),
				operation: OpCreate,
			}
			targetFile, err := os.Open(file.Target())
			if err == nil {
				content, err := ioutil.ReadAll(targetFile)
				if err != nil {
					return err
				}
				_ = targetFile.Close()
				diff.before = string(content)
				diff.operation = OpUpdate
			}

			after, err := file.RenderTemplate()
			if err != nil {
				return err
			}
			diff.after = after
			if diff.before != diff.after {
				logger.Log("Changes are required for %s!", diff.Path())
				fileDiffs = append(fileDiffs, &diff)
			}
			return nil
		}(); err != nil {
			return nil, err
		}
	}

	// run scripts.update_required and scripts.install_required to see if update is needed
	switch {
	case m.RequiresInstall():
		logger.Log("Installation is required!")
		moduleDiffs = append(moduleDiffs, &moduleDiff{
			module:    m,
			before:    StateUninstalled,
			after:     StateInstalled,
			fileDiffs: fileDiffs,
		})
	case m.RequiresUpdate():
		logger.Log("Update is required!")
		moduleDiffs = append(moduleDiffs, &moduleDiff{
			module:    m,
			before:    StateInstalled,
			after:     StateUpdated,
			fileDiffs: fileDiffs,
		})
	case len(fileDiffs) > 0:
		logger.Log("Module has file(s) needing updates!")
		moduleDiffs = append(moduleDiffs, &moduleDiff{
			module:    m,
			before:    StateInstalled,
			after:     StateInstalled,
			fileDiffs: fileDiffs,
		})
	}

	uniqueChildren := make(map[string]struct{})
	for _, mod := range m.Children() {
		if _, exists := uniqueChildren[mod.Name()]; exists {
			return nil, fmt.Errorf("error in module '%s': multiple modules defined with the same name ('%s')", m.Name(), mod.Name())
		}
		uniqueChildren[mod.Name()] = struct{}{}
		m, err := Diff(mod)
		if err != nil {
			return nil, err
		}
		moduleDiffs = append(moduleDiffs, m...)
	}

	var combinedDiffs []FileDiff
	for _, m := range moduleDiffs {
		combinedDiffs = append(combinedDiffs, m.Files()...)
	}

	filenames := make(map[string]string)
	for _, diff := range combinedDiffs {
		if existing, ok := filenames[diff.Path()]; ok {
			return nil, fmt.Errorf(
				"file '%s' must only be managed by a single module, but it is managed by both '%s' and '%s'",
				diff.Path(),
				existing,
				diff.Module().Name(),
			)
		}
		filenames[diff.Path()] = diff.Module().Name()
	}

	return moduleDiffs, nil
}
