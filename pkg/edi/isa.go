package edi

import (
	"fmt"
	"strconv"
	"strings"
)

// ISA represents the ISA EDI segment
type ISA struct {
	AuthorizationInformationQualifier string
	AuthorizationInformation          string
	SecurityInformationQualifier      string
	SecurityInformation               string
	InterchangeSenderIDQualifier      string
	InterchangeSenderID               string
	InterchangeReceiverIDQualifier    string
	InterchangeReceiverID             string
	InterchangeDate                   string
	InterchangeTime                   string
	InterchangeControlStandards       string
	InterchangeControlVersionNumber   string
	InterchangeControlNumber          int
	AcknowledgementRequested          int
	UsageIndicator                    string
	ComponentElementSeparator         string
}

// String converts ISA to its X12 single line string representation
func (s *ISA) String(delimiter string) string {
	elements := []string{
		"ISA",
		s.AuthorizationInformationQualifier,
		s.AuthorizationInformation,
		s.SecurityInformationQualifier,
		s.SecurityInformation,
		s.InterchangeSenderIDQualifier,
		s.InterchangeSenderID,
		s.InterchangeReceiverIDQualifier,
		s.InterchangeReceiverID,
		s.InterchangeDate,
		s.InterchangeTime,
		s.InterchangeControlStandards,
		s.InterchangeControlVersionNumber,
		strconv.Itoa(s.InterchangeControlNumber),
		strconv.Itoa(s.AcknowledgementRequested),
		s.UsageIndicator,
		s.ComponentElementSeparator,
	}
	return strings.Join(elements, delimiter) + "\n"
}

// Parse parses an X12 string that's split into an array into the ISA struct
func (s *ISA) Parse(elements []string) error {
	expectedNumElements := 16
	if len(elements) != expectedNumElements {
		return fmt.Errorf("ISA: Wrong number of elements, expected %d, got %d", expectedNumElements, len(elements))
	}

	var err error
	s.AuthorizationInformationQualifier = elements[0]
	s.AuthorizationInformation = elements[1]
	s.SecurityInformationQualifier = elements[2]
	s.SecurityInformation = elements[3]
	s.InterchangeSenderIDQualifier = elements[4]
	s.InterchangeSenderID = elements[5]
	s.InterchangeReceiverIDQualifier = elements[6]
	s.InterchangeReceiverID = elements[7]
	s.InterchangeDate = elements[8]
	s.InterchangeTime = elements[9]
	s.InterchangeControlStandards = elements[10]
	s.InterchangeControlVersionNumber = elements[11]
	s.InterchangeControlNumber, err = strconv.Atoi(elements[12])
	if err != nil {
		return err
	}
	s.AcknowledgementRequested, err = strconv.Atoi(elements[13])
	if err != nil {
		return err
	}
	s.UsageIndicator = elements[14]
	s.ComponentElementSeparator = elements[15]

	return nil
}
