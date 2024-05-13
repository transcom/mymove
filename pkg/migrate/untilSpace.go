package migrate

import (
	"io"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func untilSpace(in *Buffer, i int, wait time.Duration) (int, string, error) {
	var line strings.Builder
	for {
		c, err := in.Index(i)
		if err != nil {
			if err == io.EOF {
				return i, line.String(), errors.Wrap(io.EOF, "found EOF when processing stdin")
			} else if err == ErrWait {
				time.Sleep(wait)
				continue
				//nolint:revive
			} else {
				return i, line.String(), errors.Wrap(err, "received unknown error ")
			}
		}
		b := false
		switch c {
		case ' ', '\t', '\r', '\n':
			b = true
		}
		if b {
			break
		}
		i++
		line.WriteByte(c)
	}
	return i, line.String(), nil
}
