package module

import (
	"github.com/liamg/peridot/internal/pkg/config"
	"github.com/liamg/peridot/internal/pkg/system"
)

func filtersMatch(filters config.Filters) bool {

	sysInfo := system.Info()

	if len(filters.Architecture) > 0 {
		var found bool
		for _, arch := range filters.Architecture {
			if arch == sysInfo.Architecture {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(filters.OperatingSystem) > 0 {
		var found bool
		for _, os := range filters.OperatingSystem {
			if os == sysInfo.OperatingSystem {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(filters.Distribution) > 0 {
		var found bool
		for _, distro := range filters.Distribution {
			if distro == sysInfo.Distribution {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
