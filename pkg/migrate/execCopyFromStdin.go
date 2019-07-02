package migrate

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

func execCopyFromStdin(in *Buffer, i int, tablename string, columns []string, tx *pop.Connection, wait time.Duration) (int, error) {

	stmt, err := prepareCopyFromStdin(tablename, columns, tx)
	if err != nil {
		return i, errors.Wrap(err, "error preparing copy from stdin statement")
	}

	i++ // increment to next byte

	i, err = eatSpace(in, i, wait)
	if err != nil {
		return i, errors.Wrap(err, "received unknown error")
	}

	for {

		line := ""
		i, line, err = untilNewLine(in, i, wait)
		if err != nil {
			return i, err
		}

		i++ // eat new line

		// Capture end-of-data
		if line == "\\." {
			break
		}

		values := lineToValues(line)
		_, err = stmt.Exec(values...)
		if err != nil {
			// For testing might want to disable sometimes
			//fmt.Println("Error copying values into table:", err, ":", line)
			//continue
			return i, errors.Wrapf(err, "error executing copy from stdin with values %q", values)
		}
	}

	_, errFlush := stmt.Exec()
	if errFlush != nil {
		return i, errors.Wrap(errFlush, "error flushing copy from stdin")
	}

	errClose := stmt.Close()
	if errClose != nil {
		return i, errors.Wrap(errClose, "error closing copy from stdin")
	}

	return i, nil
}
