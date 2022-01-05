package builtin

import (
	"fmt"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/module"
	"github.com/liamg/peridot/internal/pkg/run"
	"github.com/liamg/peridot/internal/pkg/variable"
)

func init() {

	yayBuiltin := module.NewFactory("yay").
		WithInputs([]config.Variable{
			{
				Name:     "packages",
				Required: true,
			},
		}).
		WithRequiresInstallFunc(func(vars variable.Collection) bool {
			for _, pkg := range vars.Get("packages").AsList().All() {
				if err := run.Run(fmt.Sprintf("yay -Qi %s", pkg.AsString()), "/", false, false); err != nil {
					return true
				}
			}
			return false
		}).
		WithInstallFunc(func(vars variable.Collection) error {
			if err := run.Run("yay -Syy", "/", false, true); err != nil {
				return fmt.Errorf("failed to sync package db: %s", err)
			}
			for _, pkg := range vars.Get("packages").AsList().All() {
				if err := run.Run(fmt.Sprintf("yay -Qi %s || yay -S --noconfirm %s", pkg.AsString(), pkg.AsString()), "/", false, true); err != nil {
					return err
				}
			}
			return nil
		}).
		Build()

	module.RegisterBuiltin("yay", yayBuiltin)
}
