package config

import (
	"os"

	"github.com/liamg/peridot/internal/pkg/system"
	"github.com/liamg/peridot/internal/pkg/variable"
)

func BaseVariables() variable.Collection {

	configDir, _ := configRoot()
	homeDir, _ := os.UserHomeDir()

	info := system.Info()

	return variable.NewCollection(map[string]interface{}{
		"user_home_dir":   homeDir,
		"user_config_dir": configDir,
		"sys_os":          info.OperatingSystem,
		"sys_distro":      info.Distribution,
		"sys_arch":        info.Architecture,
	})
}
