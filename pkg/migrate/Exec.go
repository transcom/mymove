package migrate

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

func Exec(inputReader io.Reader, tx *pop.Connection, wait time.Duration) error {

	in := NewBuffer()

	go ReadInSQL(inputReader, in, true, true, true) // read in lines as a separate thread

	quoted := 0
	blocks := NewStack()
	var stmt strings.Builder

	i := 0
	for {

		//fmt.Println("Blocks:", blocks)

		char, err := in.Index(i)
		if err != nil {
			if err == io.EOF {
				//fmt.Fprintln(os.Stderr, "received EOF")
				break
			} else if err == ErrWait {
				time.Sleep(wait)
				continue
			} else {
				return errors.Wrap(err, "received unknown error")
			}
		}

		// if statement is empty, don't prefix with spaces
		if stmt.Len() == 0 && byteIsSpace(char) {
			i++
			continue
		}

		// If quoted, then see if we can unquote.
		if quoted > 0 {
			if char == '\'' {
				if in.Closed() && i+1 == in.Len() {
					quoted--
				} else {
					next, errNext := in.Index(i + 1)
					if errNext != nil {
						if errNext == io.EOF {
							break
						}
						if errNext == ErrWait {
							time.Sleep(wait)
							continue
						}
						return errors.Wrap(errNext, "received unknown error")
					}
					if next != '\'' {
						if i == 0 {
							quoted--
						} else if prev, errPrev := in.Index(i - 1); errPrev == nil && prev != '\'' {
							quoted--
						}
					}
				}
			}
			stmt.WriteByte(char)
			i++ // eat 1 character
			continue
		}

		// everything below this is unquoted

		// If not in block not quoted and on semicolon, then split statement.
		if char == ';' && blocks.Empty() && stmt.Len() > 0 {

			stmt.WriteByte(char) // append semicolon to statment

			stmtString := stmt.String()

			//fmt.Fprintln(os.Stderr, "stmt:", stmt.String())
			match := copyStdinPattern.FindStringSubmatch(stmtString)
			//fmt.Fprintln(os.Stderr, "match:", match)
			if match != nil {
				// See test to understand regex
				var errCopyFromStdin error
				i, errCopyFromStdin = execCopyFromStdin(in, i, match[4], parseColumns(match[6]), tx, wait)
				if errCopyFromStdin != nil {
					return errors.Wrap(errCopyFromStdin, "error copying from stdin")
				}
				stmt.Reset()
				continue
			}

			// If not copy from stdin
			errExec := tx.RawQuery(stmtString).Exec()
			if errExec != nil {
				return errors.Wrapf(errExec, "error executing statement: %q", stmtString)
			}
			stmt.Reset()
			i++ // forward to next character
			continue
		}

		// If not quoted and there's a quote.
		if char == '\'' {
			str, err := in.Range(i, i+2)
			if err != nil {
				if err == ErrWait {
					time.Sleep(wait)
					continue
				}
				if err != io.EOF {
					return errors.Wrap(err, "received unknown error")
				}
			}
			// err is nil
			if str == "''" {
				stmt.WriteString("''") // add the next quote too
				i += 2                 // skip forward
				continue
			}
			stmt.WriteByte(char)
			quoted++
			i++ // eat 1 character
			continue
		}

		if isAfterSpace(in, i) {
			str, err := in.Range(i, i+3)
			if err != nil {
				if err == ErrWait {
					time.Sleep(wait)
					continue
				}
				if err != io.EOF {
					return errors.Wrap(err, "received unknown error")
				}
			}
			if err == nil && byteIsSpace(str[2]) {
				pushBlock := false
				if str[0:2] == "DO" {
					if !hasPrefix(stmt.String(), "INSERT INTO") {
						pushBlock = true
					}
				} else if str[0:2] == "AS" {
					if hasPrefix(stmt.String(), "CREATE OR REPLACE FUNCTION") || hasPrefix(stmt.String(), "CREATE FUNCTION") {
						pushBlock = true
					}
				}

				if pushBlock {
					stmt.WriteString(str[0:2])
					i += 2

					i, err = eatSpace(in, i, wait)
					if err != nil {
						return errors.Wrap(err, "received unknown error")
					}

					stmt.WriteByte(' ') // add just 1 space

					block := ""
					i, block, err = untilSpace(in, i, wait)
					if err != nil {
						return errors.Wrap(err, "error reading block")
					}
					stmt.WriteString(block) // add token
					stmt.WriteRune('\n')    // add trailing new line
					blocks.Push(block)
					//fmt.Println("Blocks:", blocks)
					continue
				}
			}
		}

		// let's see if we match the last block
		if !blocks.Empty() {
			lastBlock := blocks.Last()
			str, err := in.Range(i, i+len(lastBlock))
			if err != nil {
				if err == ErrWait {
					time.Sleep(wait)
					continue
				}
				// if there is EOF and there are still blocks, then that's an issue
				return errors.Wrap(err, fmt.Sprintf("received unknown error with blocks %q left to process", blocks.Slice()))
			}
			if strings.ToUpper(str) == strings.ToUpper(lastBlock) {
				i += len(lastBlock)
				stmt.WriteString(lastBlock)
				blocks.Pop()
				continue
			}
		}

		// if nothing special simply add the character and increment the cursor
		stmt.WriteByte(char)
		i++

	}

	// If final statement did not terminate in semicolon.
	if stmt.Len() > 0 {
		if str := strings.TrimSpace(stmt.String()); len(str) > 0 {
			errExecFinalStmt := tx.RawQuery(str).Exec()
			if errExecFinalStmt != nil {
				return errors.Wrapf(errExecFinalStmt, "error executing final statement: %q", str)
			}
			stmt.Reset()
		}
	}
	return nil
}
