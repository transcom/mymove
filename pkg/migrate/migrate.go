package migrate

import (
	"regexp"
)

var (
	copyStdinPattern = regexp.MustCompile("^(\\s*)(COPY)(\\s+)([a-zA-Z._]+)(\\s*)\\((.+)\\)(\\s+)(FROM)(\\s+)(stdin)(\\s*)(;)(\\s*)$")
)
