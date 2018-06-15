package edi

import (
	"fmt"
	"strconv"
	"strings"
)

// LX represents the LX EDI segment
type LX struct {
	AssignedNumber int
}

// String converts LX to its X12 single line string representation
func (s *LX) String(delimiter string) string {
	return strings.Join([]string{"LX", strconv.Itoa(s.AssignedNumber)}, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the LX struct
func (s *LX) Parse(parts []string) error {
	expectedNumElements := 1
	if len(parts) != expectedNumElements {
		return fmt.Errorf("LX: Wrong number of elements, expected %d, got %d", expectedNumElements, len(parts))
	}

	var err error
	s.AssignedNumber, err = strconv.Atoi(parts[0])
	return err
}
