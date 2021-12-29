package cmd

import (
	"github.com/liamg/peridot/internal/pkg/module"
	"github.com/liamg/tml"
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

			diffs, err := module.Diff(root)
			if err != nil {
				fail(err)
			}

			changeCount := len(diffs)

			if changeCount == 0 {
				tml.Println("<yellow><bold>Nothing to do, no changes necessary.</bold></yellow>")
				return
			}

			for _, diff := range diffs {
				diff.Print(fullContentDiffs)
			}

			if changeCount == 1 {
				tml.Printf("\n<yellow><bold>%d module has pending changes.</bold></yellow>\n", changeCount)
			} else {
				tml.Printf("\n<yellow><bold>%d modules have pending changes.</bold></yellow>\n", changeCount)
			}
		},
	}
	diffCmd.Flags().BoolVarP(&fullContentDiffs, "show-content", "s", fullContentDiffs, "Show full git-style file content diffs.")
	rootCmd.AddCommand(diffCmd)
}
