package migrate

import (
	"bufio"
	"database/sql"
	"io"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	// This denotes the end of a copy-from-stdin data
	endCopyFromStdin = "\\."
)

// Exec executes a query
func Exec(inputReader io.Reader, tx *pop.Connection, wait time.Duration, logger *zap.Logger) error {

	scanner := bufio.NewScanner(inputReader)
	lines := make(chan string, 1000)
	dropComments := true
	dropSearchPath := true

	if scanner.Scan() {
		line := scanner.Text()
		if 0 == strings.Index(line, "-- POP RAW MIGRATION --") {
			var sb strings.Builder
			for scanner.Scan() {
				sb.WriteString(scanner.Text())
				sb.WriteString("\n")
			}
			_, err := tx.Store.Exec(sb.String())
			return err
		}
		lines <- ReadInSQLLine(line, dropComments, dropSearchPath)
	}

	// read values out of the buffer
	go func() {
		for scanner.Scan() {
			lines <- ReadInSQLLine(scanner.Text(), dropComments, dropSearchPath)
		}
		close(lines)
	}()

	statements := make(chan string, 1000)
	go SplitStatements(lines, statements, wait, logger)

	var match []string
	var preparedStmt *sql.Stmt
	for stmt := range statements {

		//if it is COPY statement then assume rest of the statements are data rows
		if match == nil {
			match = copyStdinPattern.FindStringSubmatch(stmt)
		}

		// COPY statements logic
		if match != nil {

			// Capture end of copy-from-stdin data and leave this loop
			if stmt == endCopyFromStdin {

				// Flush the statment to ensure nothing is still being buffered
				_, errFlush := preparedStmt.Exec()
				if errFlush != nil {
					return errors.Wrap(errFlush, "error flushing copy from stdin")
				}

				// Manually close the statement after flushing
				errClose := preparedStmt.Close()
				if errClose != nil {
					return errors.Wrap(errClose, "error closing copy from stdin")
				}

				// Leave the COPY statement logic
				match = nil
				preparedStmt = nil
				continue
			}

			if preparedStmt == nil {
				// prepare statement so we can insert data rows
				var errPreparedStmt error
				preparedStmt, errPreparedStmt = prepareCopyFromStdin(match[4], parseColumns(match[6]), tx)
				if errPreparedStmt != nil {
					return errors.Wrap(errPreparedStmt, "error preparing copy from stdin statement")
				}
			} else {
				values := lineToValues(stmt)
				_, errPreparedStmtExec := preparedStmt.Exec(values...)
				if errPreparedStmtExec != nil {
					return errors.Wrapf(errPreparedStmtExec, "error executing copy from stdin with values %q", values)
				}
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
