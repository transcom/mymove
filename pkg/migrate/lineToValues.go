package migrate

import (
	"strconv"
	"strings"
)

func lineToValues(line string) []interface{} {
	parts := strings.Split(line, "\t")
	values := make([]interface{}, 0, len(parts))
	for _, s := range parts {
		if s == "\\N" {
			values = append(values, nil)
		} else if n, err := strconv.Atoi(s); err == nil {
			values = append(values, int64(n))
		} else if f, err := strconv.ParseFloat(s, 64); err == nil {
			values = append(values, f)
		} else {
			values = append(values, s)
		}
	}
	return values
}
