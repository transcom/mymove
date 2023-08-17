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
	elementsLength := len(elements)
	expectedMinimumNumElements := 3 // ElementPositionInSegment might be blank

	if elementsLength < expectedMinimumNumElements {
		return fmt.Errorf("AK4: Wrong number of fields, expected at least %d, got %d %q %v", expectedMinimumNumElements, len(elements), elements, s)
	}

	var err error
	s.PositionInSegment, err = strconv.Atoi(elements[0])
	if err != nil {
		return err
	}

	if elements[1] == "" {
		s.ElementPositionInSegment = 0
	} else {
		s.ElementPositionInSegment, err = strconv.Atoi(elements[1])
		if err != nil {
			return err
		}
	}

	if elementsLength == 3 {
		// If we only have 3 segments then it must be the required DataElementSyntaxErrorCode
		s.DataElementSyntaxErrorCode = elements[2]
	} else if elementsLength == 4 {
		s.DataElementReferenceNumber, err = strconv.Atoi(elements[2])
		if err != nil {
			return err
		}
		s.CopyOfBadDataElement = elements[3]
	} else if elementsLength == 5 {
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

		s.DataElementSyntaxErrorCode = elements[3]
		// Not sure here what field to fill but the validation is more permissive
		s.CopyOfBadDataElement = elements[4]

	} else {
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
	}

	return nil
}
