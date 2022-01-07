package cmd

import (
	"fmt"
	"os"

	_ "github.com/liamg/peridot/internal/pkg/builtins"
	"github.com/liamg/peridot/internal/pkg/log"
	"github.com/liamg/tml"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "peridot",
	Short: "Manage dotfiles and user environments across machines, OSes, users and more.",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		if disableANSI {
			tml.DisableFormatting()
		}
		if debugMode {
			log.Enable()
		}
	},
}

var disableANSI bool
var debugMode bool

func init() {
	rootCmd.PersistentFlags().BoolVar(&disableANSI, "no-ansi", disableANSI, "Disable ANSI colour/formatting codes in output.")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", debugMode, "Enable debug output.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func fail(reason interface{}) {
	fmt.Fprintln(os.Stderr, tml.Sprintf("<red><bold>Error: %s</bold></red>", reason))
	os.Exit(1)
}
