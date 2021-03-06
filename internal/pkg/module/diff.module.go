package module

import (
	"fmt"

	"github.com/liamg/peridot/internal/pkg/log"
	"github.com/liamg/tml"
)

type moduleDiff struct {
	module    Module
	fileDiffs []FileDiff
	before    State
	after     State
}

func (d *moduleDiff) Files() []FileDiff {
	return d.fileDiffs
}

func (d *moduleDiff) Module() Module {
	return d.module
}

func (d *moduleDiff) Before() State {
	return d.before
}

func (d *moduleDiff) After() State {
	return d.after
}

func (d *moduleDiff) Print(withContent bool) {
	for _, f := range d.fileDiffs {
		f.Print(withContent)
	}
	if d.before != d.after {
		switch d.after {
		case StateInstalled:
			tml.Printf("<green>[Module %s] Requires install.</green>\n", d.module.Name())
		case StateUninstalled:
			tml.Printf("<red>[Module %s] Requires uninstall.</red>\n", d.module.Name())
		case StateUpdated:
			tml.Printf("<yellow>[Module %s] Requires updated.</yellow>\n", d.module.Name())
		}
	}

}

func (d *moduleDiff) Apply() error {
	logger := log.NewLogger(d.module.Name())
	for _, f := range d.fileDiffs {
		if err := f.Apply(); err != nil {
			return err
		}
	}
	if d.after != d.before {
		switch d.after {
		case StateInstalled:
			logger.Log("Installing module...")
			if err := d.module.Install(); err != nil {
				return err
			}
		case StateUpdated:
			logger.Log("Updating module...")
			if err := d.module.Update(); err != nil {
				return err
			}
		case StateUninstalled:
			logger.Log("Uninstalling module...")
			return fmt.Errorf("uninstallation is currently not supported")
		default:
			return fmt.Errorf("cannot support state 0x%X for module %s", d.after, d.module.Name())
		}
	}
	if len(d.fileDiffs) > 0 {
		logger.Log("Finalising module after file changes...")
		return d.module.AfterFileChange()
	}
	logger.Log("Application complete!...")
	return nil
}
