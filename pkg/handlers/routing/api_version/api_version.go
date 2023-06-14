package apiversion

import (
	"strings"
)

type Flag byte

const (
	NoneSpecified Flag = iota
	PrimeVersion1
	PrimeVersion2
)

func DetermineAPIVersion(path string) Flag {
	if strings.Contains(path, "/prime/v1/") {
		return PrimeVersion1
	}
	if strings.Contains(path, "/prime/v2/") {
		return PrimeVersion2
	}
	return NoneSpecified
}
