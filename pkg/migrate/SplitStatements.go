package migrate

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

import (
	"github.com/pkg/errors"
)

// SplitStatements splits a string of SQL statements into a slice with each element being a statement.
func SplitStatements(lines chan string, statements chan string) {

	in := NewBuffer()

	go func() {
		for line := range lines {
			in.WriteString(line + "\n")
		}
		in.Close()
	}()

	//statements := make([]string, 0)
	quoted := 0
	blocks := NewStack()
	var stmt strings.Builder
	i := 0
	for {
		// Get previous and current characters
		c, err := in.Index(i)
		if err != nil {
			if err == io.EOF {
				//fmt.Fprintln(os.Stderr, "received EOF")
				break
			} else if err == ErrWait {
				fmt.Fprintln(os.Stderr, "waiting for 10 milliseconds")
				time.Sleep(time.Millisecond * 10)
				continue
			} else {
				fmt.Fprintln(os.Stderr, errors.Wrap(err, "received unknown error "))
				close(statements)
				return
			}
		}

		// if statement is empty, don't prefix with spaces
		if stmt.Len() == 0 && byteIsSpace(c) {
			i++
			continue
		}

		// If not in block not quoted and on semicolon, then split statement.
		if blocks.Empty() && quoted == 0 && c == ';' {
			str := strings.TrimSpace(stmt.String() + ";")
			if len(str) > 0 {
				//statements = append(statements, stmt)
				statements <- str
				stmt.Reset()
			}
			i++ // eat 1 character
			continue
		}

		// If quoted, then see if we can unquote.
		if quoted > 0 {
			stmt.WriteByte(c)
			if c == '\'' {
				if in.Closed() && i+1 == in.Len() {
					quoted--
				} else {
					next, err := in.Index(i + 1)
					if err != nil {
						if err == io.EOF {
							break
						}
						if err == ErrWait {
							fmt.Fprintln(os.Stderr, "waiting for 10 milliseconds")
							time.Sleep(time.Millisecond * 10)
							continue
						}
						fmt.Fprintln(os.Stderr, errors.Wrap(err, "received unknown error"))
						close(statements)
						return
					}
					if next != '\'' {
						if i == 0 {
							quoted--
						} else if prev, err := in.Index(i - 1); err == nil && prev != '\'' {
							quoted--
						}
					}
				}
			}
			i++ // eat 1 character
			continue
		}

		// If not quoted and there's a quote.
		if c == '\'' {
			str, err := in.Range(i, i+2)
			if err != nil {
				if err == ErrWait {

					time.Sleep(time.Millisecond * 10)
					continue
				}
				if err != io.EOF {
					fmt.Fprint(os.Stderr, errors.Wrap(err, "received unknown error"))
					close(statements)
					return
				}
			}
			// err is nil
			if str == "''" {
				stmt.WriteString("''") // add the next quote too
				i += 2                 // skip forward
				continue
			}
			stmt.WriteByte(c)
			quoted++
			i++ // eat 1 character
			continue
		}

		if isAfterSpace(in, i) {
			str, err := in.Range(i, i+3)
			if err != nil {
				if err == ErrWait {
					fmt.Fprintln(os.Stderr, "waiting for 10 milliseconds")
					time.Sleep(time.Millisecond * 10)
					continue
				}
				if err != io.EOF {
					fmt.Fprint(os.Stderr, errors.Wrap(err, "received unknown error"))
					close(statements)
					return
				}
			}
			if err == nil && byteIsSpace(str[2]) {
				if str[0:2] == "DO" || (str[0:2] == "AS" && (hasPrefix(stmt.String(), "CREATE OR REPLACE FUNCTION") || hasPrefix(stmt.String(), "CREATE FUNCTION"))) {
					stmt.WriteString(str[0:2])
					i += 2
					for {
						c2, errIndex := in.Index(i)
						if errIndex != nil {
							if errIndex == ErrWait {
								fmt.Fprintln(os.Stderr, "waiting for 10 milliseconds")
								time.Sleep(time.Millisecond * 10)
								continue
							} else {
								fmt.Fprint(os.Stderr, errors.Wrap(errIndex, "received unknown error "))
								close(statements)
								return
							}
						}
						b := true
						switch c2 {
						case ' ', '\t', '\r', '\n':
							b = false
						}
						if b {
							break
						}
						i++
					}
					stmt.WriteString(" ") // add just 1 space
					block := ""
					for {
						c3, err := in.Index(i)
						if err != nil {
							if err == ErrWait {
								fmt.Fprintln(os.Stderr, "waiting for 10 milliseconds")
								time.Sleep(time.Millisecond * 10)
								continue
							} else {
								fmt.Fprint(os.Stderr, errors.Wrap(err, "received unknown error "))
								close(statements)
								return
							}
						}
						b := false
						switch c3 {
						case ' ', '\t', '\r', '\n':
							b = true
						default:
							block += string(c3)
						}
						if b {
							break
						}
						i++
					}

					stmt.WriteString(block) // add token
					stmt.WriteRune('\n')    // add trailing new line
					blocks.Push(strings.TrimSpace(block))
					continue
				}
			}
		}

		// let's see if we match the last block
		if !blocks.Empty() {
			lastBlock := strings.ToUpper(blocks.Last())
			str, err := in.Range(i, i+len(lastBlock))
			if err != nil {
				if err == ErrWait {
					fmt.Fprintln(os.Stderr, "waiting for 10 milliseconds")
					time.Sleep(time.Millisecond * 10)
					continue
				}
				// if there is EOF and there are still blocks, then that's an issue
				fmt.Fprint(os.Stderr, errors.Wrap(err, fmt.Sprintf("received unknown error with blocks %q left to process", blocks)))
				close(statements)
				return
			}
			if strings.ToUpper(str) == strings.ToUpper(lastBlock) {
				i += len(lastBlock)
				stmt.WriteString(blocks.Last())
				blocks.Pop()
				continue
			}
		}

		// if nothing special simply add the character and increment the cursor
		stmt.WriteByte(c)
		i++
	}

	// If final statement did not terminate in semicolon.
	if stmt.Len() > 0 {
		str := strings.TrimSpace(stmt.String())
		if len(str) > 0 {
			//statements = append(statements, stmt)
			statements <- str
			stmt.Reset()
		}
	}
	close(statements)
	//return statements
}
