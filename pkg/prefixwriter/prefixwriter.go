// https://github.com/kvz/logstreamer/blob/master/logstreamer.go
package prefixwriter

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

var _ io.Writer = &PrefixWriter{}

type PrefixWriter struct {
	writer io.Writer
	buf    *bytes.Buffer
	prefix string
}

func New(writer io.Writer, prefix string) *PrefixWriter {
	streamer := &PrefixWriter{
		writer: writer,
		prefix: prefix,
		buf:    bytes.NewBuffer([]byte("")),
	}

	return streamer
}

func (l *PrefixWriter) Write(p []byte) (n int, err error) {
	if n, err = l.buf.Write(p); err != nil {
		return
	}

	err = l.writeLines()
	return
}

func (l *PrefixWriter) Close() error {
	if err := l.flush(); err != nil {
		return err
	}
	l.buf = bytes.NewBuffer([]byte(""))
	return nil
}

func (l *PrefixWriter) flush() error {
	p := make([]byte, l.buf.Len())
	if _, err := l.buf.Read(p); err != nil {
		return err
	}

	line := string(p)
	if len(line) > 0 {
		if !strings.HasSuffix(line, "\n") {
			line += "\n"
		}

		l.out(line)
	}
	return nil
}

func (l *PrefixWriter) writeLines() error {
	for {
		line, err := l.buf.ReadString('\n')

		if len(line) > 0 {
			if strings.HasSuffix(line, "\n") {
				l.out(line)
			} else {
				// put back into buffer, it's not a complete line yet
				//  Close() or Flush() have to be used to flush out
				//  the last remaining line if it does not end with a newline
				if _, err := l.buf.WriteString(line); err != nil {
					return err
				}
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (l *PrefixWriter) out(str string) {
	if len(str) < 1 {
		return
	}

	fmt.Fprintf(l.writer, "%s %s", l.prefix, str)
}
