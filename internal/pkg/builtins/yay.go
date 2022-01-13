package builtin

import (
	"fmt"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/module"
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
		WithRequiresInstallFunc(func(r *module.Runner, vars variable.Collection) (bool, error) {
			for _, pkg := range vars.Get("packages").AsList().All() {
				if err := r.Run(fmt.Sprintf("yay -Qi %s > /dev/null", pkg.AsString()), false); err != nil {
					return true, nil
				}
			}
			return false, nil
		}).
		WithInstallFunc(func(r *module.Runner, vars variable.Collection) error {
			if err := r.Run("yay -Syy", false); err != nil {
				return fmt.Errorf("failed to sync package db: %w", err)
			}
			for _, pkg := range vars.Get("packages").AsList().All() {
				if err := r.Run(fmt.Sprintf("yay -Qi %s >/dev/null", pkg.AsString()), false); err != nil {
					if err := r.Run(fmt.Sprintf("yay -S --noconfirm %s", pkg.AsString()), true); err != nil {
						return err
					}
				}
			}
			return nil
		}).
		Build()

	module.RegisterBuiltin("yay", yayBuiltin)
}
