package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const (
	// VersionTimeFormat is the Go time format for creating a version number.
	VersionTimeFormat string = "20060102150405"
)

// checks if string is null, empty, or whitespace
func IsNullOrWhiteSpace(s *string) bool {
	return s == nil || len(strings.TrimSpace(*s)) == 0
}

func AppendTimestampToFilename(fileName string) string {
	now := time.Now()
	timestamp := now.Format(VersionTimeFormat) // ISO Format
	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)
	newFileName := fmt.Sprintf("%s-%s%s", name, timestamp, ext)
	return newFileName
}
