package cmd

import (
	"fmt"
	"os"

	"github.com/liamg/peridot/internal/pkg/module"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "diff",
		Short: "Compare the desired state as dictated by your peridot templates and config files with the actual local state.",
		Run: func(cmd *cobra.Command, args []string) {
			root, conf, err := module.ParseRoot()
			if err != nil {
				fail(err)
			}
			// TODO: show diff
			_ = root
			_ = conf
		},
	})
}

func fail(reason interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", reason)
	os.Exit(1)
}
