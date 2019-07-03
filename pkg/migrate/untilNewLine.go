package migrate

import (
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func untilNewLine(in *Buffer, i int, wait time.Duration) (int, string, error) {
	var line strings.Builder
	for {
		c, err := in.Index(i)
		if err != nil {
			if err == io.EOF {
				return i, line.String(), errors.Wrap(io.EOF, "found EOF when processing stdin")
			} else if err == ErrWait {
				time.Sleep(wait)
				continue
			} else {
				return i, line.String(), errors.Wrap(err, "received unknown error ")
			}
		}
		if c == '\n' {
			break
		}
		line.WriteByte(c)
		i++
	}
	return i, line.String(), nil
}
