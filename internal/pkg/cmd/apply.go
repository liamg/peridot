package cmd

import (
	"fmt"

	"github.com/liamg/peridot/internal/pkg/module"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Compare the desired state as dictated by your peridot templates and config files with the actual local state.",
		Run: func(cmd *cobra.Command, args []string) {
			root, err := module.LoadRoot()
			if err != nil {
				fail(err)
			}

			moduleDiffs, fileDiffs, err := root.Diff()
			if err != nil {
				fail(err)
			}

			changeCount := len(moduleDiffs) + len(fileDiffs)

			if changeCount == 0 {
				fmt.Println("Nothing to do, no changes necessary.")
				return
			}

			for _, fileDiff := range fileDiffs {
				if err := fileDiff.Apply(); err != nil {
					fail(err)
				}
			}

			for _, moduleDiff := range moduleDiffs {
				if err := moduleDiff.Apply(); err != nil {
					fail(err)
				}
			}

			fmt.Printf("\n%d changes applied successfully.", changeCount)
		},
	}
	rootCmd.AddCommand(applyCmd)
}
