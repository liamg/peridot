package validate

import (
	"fmt"

	"github.com/liamg/peridot/internal/pkg/variable"
)

type Validator struct {
	vars variable.Collection
}

func New(vars variable.Collection) *Validator {
	return &Validator{
		vars: vars,
	}
}

func (v *Validator) EnsureDefined(names ...string) error {
	for _, name := range names {
		if !v.vars.Has(name) {
			return fmt.Errorf("required variable '%s' is not defined", name)
		}
	}
	return nil
}
