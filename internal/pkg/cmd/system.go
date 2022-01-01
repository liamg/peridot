package cmd

import (
	"github.com/liamg/peridot/internal/pkg/system"
	"github.com/liamg/tml"
	"github.com/spf13/cobra"
)

func init() {
	systemCmd := &cobra.Command{
		Use:   "system",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			info := system.Info()
			tml.Printf("Architecture:     <bold>%s</bold>\n", info.Architecture)
			tml.Printf("Operating System: <bold>%s</bold>\n", info.OperatingSystem)
			distro := info.Distribution
			if distro == "" {
				distro = "n/a"
			}
			tml.Printf("Distribution:     <bold>%s</bold>\n", distro)
		},
	}
	rootCmd.AddCommand(systemCmd)
}
