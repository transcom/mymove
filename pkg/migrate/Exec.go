package migrate

import (
	"bufio"
	"io"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

func Exec(inputReader io.Reader, tx *pop.Connection, wait time.Duration) error {

	in := NewBuffer()

	go ReadInSQL(inputReader, in, true, true, true) // read in lines as a separate thread

	lines := make(chan string, 1000)
	// read values out of the buffer
	go func() {
		formattedSQL := in.String()
		scanner := bufio.NewScanner(strings.NewReader(formattedSQL))
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	statements := make(chan string, 1000)
	go SplitStatements(lines, statements)
	counter := 0
	for stmt := range statements {
		//if it is COPY statement then assume rest of the statements are part of copy and execute then as part of stdin
		match := copyStdinPattern.FindStringSubmatch(stmt)

		if match != nil {
			//create buffer of remaining statements
			var copyStmts Buffer

			//range over the channel again
			for copyStmt := range statements {
				isCopy := copyStdinPattern.FindStringSubmatch(copyStmt)

				//skip all non copy statements
				if isCopy == nil {
					continue
				} else {
					_, err := copyStmts.WriteString(copyStmt)
					if err != nil {
						return errors.Wrap(err, "error copying from stdin")
					}
				}
			}

			// See test to understand regex
			var errCopyFromStdin error
			_, errCopyFromStdin = execCopyFromStdin(&copyStmts, counter, match[4], parseColumns(match[6]), tx, wait)
			if errCopyFromStdin != nil {
				return errors.Wrap(errCopyFromStdin, "error copying from stdin")
			}
			//exit loop after encountering the first COPY statement
			break
		}

		//if it is a regular statement then execute it as raw sql
		errExec := tx.RawQuery(stmt).Exec()
		if errExec != nil {
			return errors.Wrapf(errExec, "error executing statement: %q", stmt)
		}
		counter++
	}
	return nil
}