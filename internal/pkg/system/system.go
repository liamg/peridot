package system

import (
	"os"
	"runtime"
	"strings"
)

type systemInfo struct {
	OperatingSystem string
	Distribution    string
	Architecture    string
}

var systemInfoCache *systemInfo

func Info() systemInfo {
	if systemInfoCache == nil {
		systemInfoCache = &systemInfo{
			OperatingSystem: runtime.GOOS,
			Architecture:    runtime.GOARCH,
			Distribution:    getDistro(),
		}
	}
	return *systemInfoCache
}

func getDistro() string {
	for _, file := range []string{
		"/etc/os-release",
		"/usr/lib/os-release",
		"/etc/initrd-release",
	} {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "ID=") {
				return strings.Trim(line[3:], `"'`)
			}
		}
	}
	return ""
}
