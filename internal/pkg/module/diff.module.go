package module

import "fmt"

type moduleDiff struct {
	module Module
	before State
	after  State
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

func (d *moduleDiff) Print() {
	switch {
	case d.before == d.after:
		fmt.Printf("Module %s has no known changes.\n", d.module.Name())
	case d.after == StateInstalled:
		fmt.Printf("Module %s will be installed.\n", d.module.Name())
	case d.after == StateUninstalled:
		fmt.Printf("Module %s will be uninstalled.\n", d.module.Name())
	case d.after == StateUpdated:
		fmt.Printf("Module %s will be updated.\n", d.module.Name())
	}

}

func (d *moduleDiff) Apply() error {
	switch d.after {
	case StateInstalled:
		return d.module.Install()
	case StateUpdated:
		return d.module.Update()
	case StateUninstalled:
		return fmt.Errorf("uninstallation is currently not supported")
	default:
		return fmt.Errorf("cannot support state 0x%X for module %s", d.after, d.module.Name())
	}
}
