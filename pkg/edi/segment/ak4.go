package edisegment

import (
	"fmt"
	"strconv"
)

// AK4 represents the AK4 EDI segment
type AK4 struct {
	PositionInSegment                       int    `validate:"min=1,max=99"`
	ElementPositionInSegment                int    `validate:"min=1,max=99"`
	ComponentDataElementPositionInComposite int    `validate:"omitempty,min=1,max=99"`
	DataElementReferenceNumber              int    `validate:"omitempty,min=1,max=9999"`
	DataElementSyntaxErrorCode              string `validate:"min=1,max=3"`
	CopyOfBadDataElement                    string `validate:"omitempty,max=99"`
}

// StringArray converts AK4 to an array of strings
func (s *AK4) StringArray() []string {
	// For the optional int fields, make sure we map a zero to the empty string.
	var strComponentDataElementPositionInComposite string
	if s.ComponentDataElementPositionInComposite != 0 {
		strComponentDataElementPositionInComposite = strconv.Itoa(s.ComponentDataElementPositionInComposite)
	}

	var strDataElementReferenceNumber string
	if s.DataElementReferenceNumber != 0 {
		strDataElementReferenceNumber = strconv.Itoa(s.DataElementReferenceNumber)
	}

	return []string{
		"AK4",
		strconv.Itoa(s.PositionInSegment),
		strconv.Itoa(s.ElementPositionInSegment),
		strComponentDataElementPositionInComposite,
		strDataElementReferenceNumber,
		s.DataElementSyntaxErrorCode,
		s.CopyOfBadDataElement,
	}
}

// Parse parses an X12 string that's split into an array into the AK4 struct
func (s *AK4) Parse(elements []string) error {
	expectedNumElements := 6
	if len(elements) != expectedNumElements {
		return fmt.Errorf("AK4: Wrong number of fields, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.PositionInSegment, err = strconv.Atoi(elements[0])
	if err != nil {
		return err
	}

	s.ElementPositionInSegment, err = strconv.Atoi(elements[1])
	if err != nil {
		return err
	}

	// For the optional int fields, make sure we map an empty string to a zero.
	strComponentDataElementPositionInComposite := elements[2]
	if strComponentDataElementPositionInComposite == "" {
		s.ComponentDataElementPositionInComposite = 0
	} else {
		s.ComponentDataElementPositionInComposite, err = strconv.Atoi(strComponentDataElementPositionInComposite)
		if err != nil {
			return err
		}
	}

	strDataElementReferenceNumber := elements[3]
	if strDataElementReferenceNumber == "" {
		s.DataElementReferenceNumber = 0
	} else {
		s.DataElementReferenceNumber, err = strconv.Atoi(strDataElementReferenceNumber)
		if err != nil {
			return err
		}
	}

	s.DataElementSyntaxErrorCode = elements[4]
	s.CopyOfBadDataElement = elements[5]

	return nil
}
