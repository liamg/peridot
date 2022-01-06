package builtin

import (
	"fmt"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/module"
	"github.com/liamg/peridot/internal/pkg/variable"
)

func init() {

	pacmanBuiltin := module.NewFactory("pacman").
		WithInputs([]config.Variable{
			{
				Name:     "packages",
				Required: true,
			},
		}).
		WithRequiresInstallFunc(func(r *module.Runner, vars variable.Collection) (bool, error) {
			for _, pkg := range vars.Get("packages").AsList().All() {
				if err := r.Run(fmt.Sprintf("pacman -Qi %s >/dev/null", pkg.AsString()), false); err != nil {
					return true, nil
				}
			}
			return false, nil
		}).
		WithInstallFunc(func(r *module.Runner, vars variable.Collection) error {
			if err := r.Run("pacman -Syy", true); err != nil {
				return fmt.Errorf("failed to sync package db: %s", err)
			}
			for _, pkg := range vars.Get("packages").AsList().All() {
				if err := r.Run(fmt.Sprintf("pacman -Qi %s >/dev/null", pkg.AsString()), false); err != nil {
					if err := r.Run(fmt.Sprintf("pacman -S --noconfirm %s", pkg.AsString()), true); err != nil {
						return err
					}
				}
			}
			return nil
		}).
		Build()

	module.RegisterBuiltin("pacman", pacmanBuiltin)
}
