package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/tml"
	"github.com/spf13/cobra"
)

func init() {
	var force bool
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialises a new peridot config for the local user environment.",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			configPath, err := config.Path()
			if err != nil {
				fail(fmt.Sprintf("Cannot locate config path: %s", err))
			}
			if _, err := os.Stat(configPath); err == nil {
				if force {
					if err := os.RemoveAll(filepath.Dir(configPath)); err != nil {
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
			tml.Printf("<green><bold>New configuration file and git repository initialised at %s</bold></green>\n", path)
		},
	}
	initCmd.Flags().BoolVarP(&force, "force", "f", force, "Force peridot to overwrite an existing config and reinitialise a fresh one.")
	rootCmd.AddCommand(initCmd)
}
