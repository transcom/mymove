package migrate

import (
	"unicode"
)

func byteIsSpace(b byte) bool {
	return unicode.IsSpace([]rune(string([]byte{b}))[0])
}
