package config

import (
	"os"

	"github.com/liamg/peridot/internal/pkg/variable"
)

func BaseVariables() variable.Collection {

	configDir, _ := configRoot()
	homeDir, _ := os.UserHomeDir()

	return variable.NewCollection(map[string]interface{}{
		"user_home_dir":   homeDir,
		"user_config_dir": configDir,
	})
}
