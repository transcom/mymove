package edisegment

import (
	"fmt"
	"strconv"
	"strings"
)

// SE represents the SE EDI segment
type SE struct {
	NumberOfIncludedSegments    int
	TransactionSetControlNumber string
}

// String converts SE to its X12 single line string representation
func (s *SE) String(delimiter string) string {
	elements := []string{
		"SE",
		strconv.Itoa(s.NumberOfIncludedSegments),
		s.TransactionSetControlNumber,
	}
	return strings.Join(elements, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the SE struct
func (s *SE) Parse(elements []string) error {
	expectedNumElements := 2
	if len(elements) != expectedNumElements {
		return fmt.Errorf("SE: Wrong number of fields, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.NumberOfIncludedSegments, err = strconv.Atoi(elements[0])
	if err != nil {
		return err
	}
	s.TransactionSetControlNumber = elements[1]
	return nil
}
