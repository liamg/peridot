package config

import (
	"os"
)

func BaseVariables() map[string]interface{} {

	configDir, _ := configRoot()
	homeDir, _ := os.UserHomeDir()

	return map[string]interface{}{
		"user_home_dir":   homeDir,
		"user_config_dir": configDir,
	}
}
