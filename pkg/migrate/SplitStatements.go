package migrate

import (
	"io"
	"strings"
	"time"
)

// SplitStatements splits a string of SQL statements into a slice with each element being a statement.
func SplitStatements(lines chan string, statements chan string, wait time.Duration) {

	in := NewBuffer()

	go func() {
		for line := range lines {
			in.WriteString(line + "\n")
		}
		in.Close()
	}()

	quoted := 0
	blocks := NewStack()
	var stmt strings.Builder

	i := 0
	var hasCopyStatment []string
	for {
		// Get previous and current characters
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

		// if statement is empty, don't prefix with spaces
		if stmt.Len() == 0 && byteIsSpace(char) {
			i++
			continue
		}

		// If not in block not quoted and on semicolon, then split statement.
		if blocks.Empty() && quoted == 0 && char == ';' {
			str := strings.TrimSpace(stmt.String() + ";")
			if len(str) > 0 {
				statements <- str
				stmt.Reset()
			}

			// check for copy statement in statements, we only need to do this once
			if hasCopyStatment == nil {
				hasCopyStatment = copyStdinPattern.FindStringSubmatch(str)
			}
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
						}
					}
				}
			}
			i++ // eat 1 character
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

		// let's see if we match the last block
		if !blocks.Empty() {
			lastBlock := strings.ToUpper(blocks.Last())
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
			if strings.ToUpper(str) == strings.ToUpper(lastBlock) {
				i += len(lastBlock)
				stmt.WriteString(blocks.Last())
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
		str := strings.TrimSpace(stmt.String())
		if len(str) > 0 {
			// if we found COPY statement in statements anywhere, we assume last rows are data rows. Split them on newline
			if hasCopyStatment != nil {
				copyStmts := strings.Split(str, "\n")
				for _, dataStmt := range copyStmts {
					// skip adding \.
					if dataStmt == "\\." {
						continue
					}
					statements <- dataStmt
				}
			} else {
				statements <- str
				stmt.Reset()
			}
		}
	}
	close(statements)
}
