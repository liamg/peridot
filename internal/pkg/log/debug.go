package log

import "fmt"

var enabled bool

func Enable() {
	enabled = true
}

func Debug(format string, args ...interface{}) {
	if !enabled {
		return
	}
	fmt.Printf(format+"\n", args)
}
