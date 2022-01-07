package log

import (
	"fmt"
	"sync"
)

type prefixedWriter struct {
	sync.Mutex
	logger *Logger
	buffer string
}

func NewPrefixedWriter(prefix string, operation string) *prefixedWriter {
	return &prefixedWriter{
		logger: NewLogger(fmt.Sprintf("%s:%s", prefix, operation)),
	}
}

func (p *prefixedWriter) Flush() {
	if p.buffer == "" {
		return
	}
	p.logger.Log("%s", p.buffer)
	p.buffer = ""
}

func (p *prefixedWriter) Write(d []byte) (n int, err error) {
	p.Lock()
	defer p.Unlock()
	for _, r := range string(d) {
		if r == '\n' {
			p.Flush()
			continue
		}
		p.buffer += string(r)
	}
	return len(d), nil
}
