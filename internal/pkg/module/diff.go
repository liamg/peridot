package module

import (
	"io/ioutil"
	"os"
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

	for _, file := range m.Files() {
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
				fileDiffs = append(fileDiffs, &diff)
			}
			return nil
		}(); err != nil {
			return nil, err
		}
	}

	// run scripts.update_required and scripts.install_required to see if update is needed
	if m.RequiresInstall() {
		moduleDiffs = append(moduleDiffs, &moduleDiff{
			module:    m,
			before:    StateUninstalled,
			after:     StateInstalled,
			fileDiffs: fileDiffs,
		})
	} else if m.RequiresUpdate() {
		moduleDiffs = append(moduleDiffs, &moduleDiff{
			module:    m,
			before:    StateInstalled,
			after:     StateUpdated,
			fileDiffs: fileDiffs,
		})
	} else if len(fileDiffs) > 0 {
		moduleDiffs = append(moduleDiffs, &moduleDiff{
			module:    m,
			before:    StateInstalled,
			after:     StateInstalled,
			fileDiffs: fileDiffs,
		})
	}

	for _, mod := range m.Children() {
		m, err := Diff(mod)
		if err != nil {
			return nil, err
		}
		moduleDiffs = append(moduleDiffs, m...)
	}

	return moduleDiffs, nil
}
