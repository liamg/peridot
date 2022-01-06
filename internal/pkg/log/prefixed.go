package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

type prefixedWriter struct {
	sync.Mutex
	prefix    string
	recipient io.Writer
	buffer    string
}

func NewDebugStreamer(prefix string) *prefixedWriter {
	if !enabled {
		return NewPrefixedWriter(prefix, ioutil.Discard)
	}
	return NewPrefixedWriter(prefix, os.Stdout)
}

func NewPrefixedWriter(prefix string, recipient io.Writer) *prefixedWriter {
	return &prefixedWriter{
		prefix:    prefix,
		recipient: recipient,
	}
}

func (p *prefixedWriter) Flush() error {
	if p.buffer == "" {
		return nil
	}
	_, err := fmt.Fprintf(p.recipient, "%s %s\n", p.prefix, p.buffer)
	p.buffer = ""
	return err
}

func (p *prefixedWriter) Write(d []byte) (n int, err error) {
	p.Lock()
	defer p.Unlock()
	var written int
	for _, r := range string(d) {
		if r == '\n' {
			if err := p.Flush(); err != nil {
				return written, err
			}
			continue
		}
		p.buffer += string(r)
	}
	return len(d), nil
}
