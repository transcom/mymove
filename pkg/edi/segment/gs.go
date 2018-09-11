package edisegment

import (
	"fmt"
	"strconv"
	"strings"
)

// GS represents the GS EDI segment
type GS struct {
	FunctionalIdentifierCode string
	ApplicationSendersCode   string
	ApplicationReceiversCode string
	Date                     string
	Time                     string
	GroupControlNumber       int
	ResponsibleAgencyCode    string
	Version                  string
}

// String converts GS to its X12 single line string representation
func (s *GS) String(delimiter string) string {
	elements := []string{
		"GS",
		s.FunctionalIdentifierCode,
		s.ApplicationSendersCode,
		s.ApplicationReceiversCode,
		s.Date,
		s.Time,
		strconv.Itoa(s.GroupControlNumber),
		s.ResponsibleAgencyCode,
		s.Version,
	}
	return strings.Join(elements, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the GS struct
func (s *GS) Parse(elements []string) error {
	expectedNumElements := 8
	if len(elements) != expectedNumElements {
		return fmt.Errorf("GS: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.FunctionalIdentifierCode = elements[0]
	s.ApplicationSendersCode = elements[1]
	s.ApplicationReceiversCode = elements[2]
	s.Date = elements[3]
	s.Time = elements[4]
	s.GroupControlNumber, err = strconv.Atoi(elements[5])
	if err != nil {
		return err
	}
	s.ResponsibleAgencyCode = elements[6]
	s.Version = elements[7]
	return nil
}
