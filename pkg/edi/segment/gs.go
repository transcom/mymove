package edisegment

import (
	"fmt"
	"strconv"
)

// GS represents the GS EDI segment
type GS struct {
	FunctionalIdentifierCode string `validate:"oneof=SI AG FA"`
	ApplicationSendersCode   string `validate:"oneof=MILMOVE 8004171844"`
	ApplicationReceiversCode string `validate:"oneof=MILMOVE 8004171844"`
	Date                     string `validate:"datetime=20060102"`
	Time                     string `validate:"datetime=1504|datetime=150405"`
	GroupControlNumber       int64  `validate:"min=1,max=999999999"`
	ResponsibleAgencyCode    string `validate:"eq=X"`
	Version                  string `validate:"eq=004010"`
}

// StringArray converts GS to an array of strings
func (s *GS) StringArray() []string {
	return []string{
		"GS",
		s.FunctionalIdentifierCode,
		s.ApplicationSendersCode,
		s.ApplicationReceiversCode,
		s.Date,
		s.Time,
		strconv.FormatInt(s.GroupControlNumber, 10),
		s.ResponsibleAgencyCode,
		s.Version,
	}
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
	s.GroupControlNumber, err = strconv.ParseInt(elements[5], 10, 64)
	if err != nil {
		return err
	}
	s.ResponsibleAgencyCode = elements[6]
	s.Version = elements[7]
	return nil
}
