package cmd

import (
	"fmt"
	"os"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/module"
	"github.com/spf13/cobra"
)

func init() {
	var force bool
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialises a new peridot config for the local user environment.",
		Run: func(cmd *cobra.Command, args []string) {
			root, err := module.LoadRoot()
			if err == nil {
				if force {
					if err := os.RemoveAll(root.Path()); err != nil {
						fail(err)
					}
				} else {
					fail("Configuration already exists. Use --force to completely remove all peridot config and create a blank config.")
				}
			} else if !os.IsNotExist(err) {
				fail(err)
			}
			path, err := config.Init()
			if err != nil {
				fail(err)
			}
			fmt.Printf("New configuration file initialised at %s\n", path)
		},
	}
	initCmd.Flags().BoolVarP(&force, "force", "f", force, "Force peridot to overwrite an existing config and reinitialise a fresh one.")
	rootCmd.AddCommand(initCmd)
}
