package migrate

import (
	"io"
	"strings"
	"time"

	"go.uber.org/zap"
)

// SplitStatements splits a string of SQL statements into a slice with each element being a statement.
func SplitStatements(lines chan string, statements chan string, wait time.Duration, logger *zap.Logger) {

	in := NewBuffer()

	go func() {
		for line := range lines {
			// Ignore empty lines when writing to the buffer
			if len(line) > 0 {

				_, err := in.WriteString(line + "\n")
				if err != nil {
					logger.Error("Failed to ignore empty lines when writing to the buffer: %s", zap.Error(err))
				}
			}
		}
		in.Close()
	}()

	quoted := 0
	quoteRunLength := 0
	blocks := NewStack()
	var stmt strings.Builder

	i := 0
	inCopyStatement := false
	for {
		// Get current character
		char, err := in.Index(i)
		if err != nil {
			if err == io.EOF {
				break
			} else if err == ErrWait {
				time.Sleep(wait)
				continue
			} else {
				close(statements)
				return
			}
		}

		if char != '\'' {
			quoteRunLength = 0
		} else {
			quoteRunLength++
		}

		// If statement is empty, don't prefix with spaces
		if stmt.Len() == 0 && byteIsSpace(char) {
			i++
			continue
		}

		// If we're in a COPY statement, see if we've reached the end stdin marker
		if inCopyStatement && char == '\n' {
			twoPrevChars, err := in.Range(i-2, i)
			if err != nil {
				if err == ErrWait {
					time.Sleep(wait)
					continue
				} else {
					close(statements)
					return
				}
			} else if twoPrevChars == "\\." {
				// Found the end stdin marker, so slice the string into lines, send to the
				// channel, and reset our copy statement boolean.
				str := strings.TrimSpace(stmt.String())
				copyStmts := strings.Split(str, "\n")
				for _, dataStmt := range copyStmts {
					statements <- dataStmt
					stmt.Reset()
				}

				inCopyStatement = false
				i++ // eat 1 character
				continue
			}
		}

		// If not in block not quoted and on semicolon, then split statement.
		if blocks.Empty() && quoted == 0 && char == ';' && !inCopyStatement {
			str := strings.TrimSpace(stmt.String() + ";")
			if len(str) > 0 { // will the len ever be zero? we're adding a semicolon?
				statements <- str
				stmt.Reset()
			}

			// Check to see if this was a COPY statement
			inCopyStatement = copyStdinPattern.MatchString(str)
			i++ // eat 1 character
			continue
		}

		// If quoted, then see if we can unquote.
		if quoted > 0 {
			stmt.WriteByte(char)
			if char == '\'' {
				if in.Closed() && i+1 == in.Len() {
					quoted--
				} else {
					next, err := in.Index(i + 1)
					if err != nil {
						if err == io.EOF {
							break
						}
						if err == ErrWait {
							time.Sleep(wait)
							continue
						}
						close(statements)
						return
					}
					if next != '\'' {
						if i == 0 {
							quoted--
						} else if prev, err := in.Index(i - 1); err == nil && prev != '\'' {
							quoted--
						} else if quoteRunLength%2 == 1 {
							// An odd number of consecutive quotes within a string literal includes an opening or closing quote
							// opening quotes are handled elsewhere, so when we get here we must have a closing quote
							quoted--
						}
					}
				}
			}
			i++ // eat 1 character
			continue
		}

		// If not quoted and there's a quote, increase our quote level.
		if char == '\'' && !inCopyStatement {
			str, err := in.Range(i, i+2)
			if err != nil {
				if err == ErrWait {
					time.Sleep(wait)
					continue
				}
				if err != io.EOF {
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
			stmt.WriteByte(char)
			quoted++
			i++ // eat 1 character
			continue
		}

		// Look for blocks of code such as "DO"
		if isAfterSpace(in, i) {
			str, err := in.Range(i, i+3)
			if err != nil {
				if err == ErrWait {
					time.Sleep(wait)
					continue
				}
				if err != io.EOF {
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
								time.Sleep(wait)
								continue
							} else {
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
								time.Sleep(wait)
								continue
							} else {
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

		// Let's see if we match the last block
		if !blocks.Empty() {
			lastBlock := blocks.Last()
			str, err := in.Range(i, i+len(lastBlock))
			if err != nil {
				if err == ErrWait {
					time.Sleep(wait)
					continue
				}
				// if there is EOF and there are still blocks, then that's an issue
				close(statements)
				return
			}
			if strings.EqualFold(str, lastBlock) {
				i += len(lastBlock)
				stmt.WriteString(blocks.Last())
				blocks.Pop()
				continue
			}
		}

		// If nothing special, simply add the character and increment the cursor
		stmt.WriteByte(char)
		i++
	}

	// If final statement did not terminate in semicolon.
	if stmt.Len() > 0 {
		str := strings.TrimSpace(stmt.String())
		if len(str) > 0 {
			statements <- str
			stmt.Reset()
		}
	}
	close(statements)
}
