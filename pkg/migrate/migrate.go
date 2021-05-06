package migrate

import (
	"regexp"
)

var (
	copyStdinPattern = regexp.MustCompile(`^(\s*)(COPY)(\s+)([a-zA-Z0-9._]+)(\s*)\((.+)\)(\s+)(FROM)(\s+)(stdin)(\s*)(;)(\s*)$`)
)
