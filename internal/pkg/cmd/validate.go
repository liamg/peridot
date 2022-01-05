package cmd

import (
	"path/filepath"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/module"
	"github.com/liamg/tml"
	"github.com/spf13/cobra"
)

func init() {
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the current configuration, taking no further action.",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			root, err := module.LoadRoot()
			if err != nil {
				fail(err)
			}
			tml.Printf("<green><bold>Configuration at '%s' appears valid.</bold></green>\n", filepath.Join(root.Path(), config.Filename))
		},
	}
	rootCmd.AddCommand(validateCmd)
}
