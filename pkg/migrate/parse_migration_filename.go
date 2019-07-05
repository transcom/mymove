package migrate

import (
	"fmt"
	"regexp"

	"github.com/gobuffalo/pop"
)

var mrx = regexp.MustCompile(`^(\d+)_([^.]+)(\.[a-z0-9]+)?(\.[a-z]+)?\.(sql|fizz)$`)

// ParseMigrationFilename parses a migration filename.
func ParseMigrationFilename(filename string) (*pop.Match, error) {

	matches := mrx.FindAllStringSubmatch(filename, -1)
	if len(matches) == 0 {
		return nil, nil
	}
	m := matches[0]

	dbType := ""
	direction := ""
	if len(m[3]) == 0 {
		dbType = "all"
		direction = "up"
	} else {
		if len(m[4]) == 0 {
			dbType = "all"
			direction = m[3][1:]
		} else {
			dbType = m[3][1:]
			if !pop.DialectSupported(dbType) {
				return nil, fmt.Errorf("unsupported dialect %s", dbType)
			}
			direction = m[4][1:]
		}
	}

	fileType := m[len(m)-1]

	if fileType == "fizz" && dbType != "all" {
		return nil, fmt.Errorf("invalid database type %q, expected \"all\" because fizz is database type independent", dbType)
	}

	match := &pop.Match{
		Version:   m[1],
		Name:      m[2],
		DBType:    dbType,
		Direction: direction,
		Type:      fileType,
	}

	return match, nil
}
