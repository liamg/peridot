package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/liamg/tml"
)

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

type Logger struct {
	prefix    string
	recipient io.Writer
}

func NewLogger(prefix string) *Logger {
	recipient := ioutil.Discard
	if enabled {
		recipient = os.Stdout
	}
	return &Logger{
		prefix:    prefix,
		recipient: recipient,
	}
}

func (l *Logger) Log(format string, args ...interface{}) {
	fmt.Fprintln(
		l.recipient,
		tml.Sprintf(
			"[<bold><yellow>%s</yellow></bold>] %s",
			l.prefix,
			fmt.Sprintf(format, args...),
		),
	)
}
