package migrate

import (
	"bufio"
	"io"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

func Exec(inputReader io.Reader, tx *pop.Connection, wait time.Duration) error {

	lines := make(chan string, 1000)
	dropComments := true
	dropSearchPath := true
	// read values out of the buffer
	go func() {
		scanner := bufio.NewScanner(inputReader)
		for scanner.Scan() {
			lines <- ReadInSQLLine(scanner.Text(), dropComments, dropSearchPath)
		}
		close(lines)
	}()

	statements := make(chan string, 1000)
	go SplitStatements(lines, statements, wait)

	var match []string
	for stmt := range statements {

		//if it is COPY statement then assume rest of the statements are data rows
		if match == nil {
			match = copyStdinPattern.FindStringSubmatch(stmt)
		}

		// COPY statements logic
		if match != nil {
			// prepare statement so we can insert data rows
			preparedStmt, err := prepareCopyFromStdin(match[4], parseColumns(match[6]), tx)
			if err != nil {
				return errors.Wrap(err, "error preparing copy from stdin statement")
			}

			values := lineToValues(stmt)
			_, err = preparedStmt.Exec(values...)
			if err != nil {
				return errors.Wrapf(err, "error executing copy from stdin with values %q", values)
			}
		} else {
			//regular statement logic executes it as raw sql
			errExec := tx.RawQuery(stmt).Exec()
			if errExec != nil {
				return errors.Wrapf(errExec, "error executing statement: %q", stmt)
			}
		}

	}
	return nil
}
