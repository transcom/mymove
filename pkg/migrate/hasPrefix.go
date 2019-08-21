package migrate

import (
	"strings"
)

func hasPrefix(str string, prefix string) bool {
	return strings.HasPrefix(strings.ToUpper(str), strings.ToUpper(prefix))
}
