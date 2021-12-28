package cmd

import (
	"fmt"
	"os"

	"github.com/liamg/peridot/internal/pkg/module"
	"github.com/spf13/cobra"
)

func init() {
	var fullContentDiffs bool
	diffCmd := &cobra.Command{
		Use:   "diff",
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
				fileDiff.Print(fullContentDiffs)
				fmt.Println("")
			}

			for _, moduleDiff := range moduleDiffs {
				moduleDiff.Print()
				fmt.Println("")
			}

			fmt.Printf("\n%d pending changes detected.", changeCount)
		},
	}
	diffCmd.Flags().BoolVarP(&fullContentDiffs, "show-content", "s", fullContentDiffs, "Show full git-style file content diffs.")
	rootCmd.AddCommand(diffCmd)
}

func fail(reason interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", reason)
	os.Exit(1)
}
