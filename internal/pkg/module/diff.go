package module

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
	Print()
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
