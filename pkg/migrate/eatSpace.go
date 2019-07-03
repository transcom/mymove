package migrate

import (
	"time"

	"github.com/pkg/errors"
)

func eatSpace(in *Buffer, i int, wait time.Duration) (int, error) {
	// eat all space until next line
	for {
		c, err := in.Index(i)
		if err != nil {
			if err == ErrWait {
				time.Sleep(wait)
				continue
			} else {
				return i, errors.Wrap(err, "received unknown error ")
			}
		}
		b := true
		switch c {
		case ' ', '\t', '\r', '\n':
			b = false
		}
		if b {
			break
		}
		i++
	}
	return i, nil
}
