package module

import (
	"fmt"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/variable"
)

func validateVariables(expectedVars []config.Variable, actual variable.Collection) error {
	for _, expected := range expectedVars {
		if expected.Required {
			if !actual.Has(expected.Name) {
				return fmt.Errorf("required variable '%s' is not defined", expected.Name)
			}
		}
	}
	return nil
}

func applyVariableDefaults(expectedVars []config.Variable, actual variable.Collection) variable.Collection {
	merged := variable.NewCollection(nil)
	for _, input := range expectedVars {
		if actual.Has(input.Name) {
			merged.Set(input.Name, actual.Get(input.Name).Interface())
		} else if !input.Required && input.Default != nil {
			merged.Set(input.Name, input.Default)
		}
	}
	merged.MergeIn(config.BaseVariables())
	return merged

}
