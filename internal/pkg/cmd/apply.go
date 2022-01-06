package cmd

import (
	"fmt"

	"github.com/liamg/peridot/internal/pkg/module"
	"github.com/liamg/tml"
	"github.com/spf13/cobra"
)

func init() {
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply all required changes as dictated by the current configuration. You can preview changes first with the 'diff' command.",
		Run: func(_ *cobra.Command, _ []string) {
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

			for _, moduleDiff := range diffs {
				tml.Printf("<yellow><bold>[Module %s] Applying changes...", moduleDiff.Module().Name())
				if err := moduleDiff.Apply(); err != nil {
					fmt.Println("")
					fail(err)
				}
				tml.Printf("\x1b[2K\r<green>[Module %s] Changes applied.\n", moduleDiff.Module().Name())
			}

			tml.Printf("\n<green><bold>%d modules applied successfully.</bold></green>\n", changeCount)
		},
	}
	rootCmd.AddCommand(applyCmd)
}
