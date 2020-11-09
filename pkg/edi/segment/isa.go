package edisegment

import (
	"fmt"
	"strconv"
)

// ISA represents the ISA EDI segment
type ISA struct {
	AuthorizationInformationQualifier string `validate:"eq=00"`
	AuthorizationInformation          string `validate:"eq=0084182369"`
	SecurityInformationQualifier      string `validate:"eq=00"`
	SecurityInformation               string `validate:"eq=0000000000"`
	InterchangeSenderIDQualifier      string `validate:"eq=ZZ"`
	InterchangeSenderID               string `validate:"eq=MILMOVE        "`
	InterchangeReceiverIDQualifier    string `validate:"eq=12"`
	InterchangeReceiverID             string `validate:"eq=8004171844     "`
	InterchangeDate                   string `validate:"datetime=060102"`
	InterchangeTime                   string `validate:"datetime=1504"`
	InterchangeControlStandards       string `validate:"eq=U"`
	InterchangeControlVersionNumber   string `validate:"eq=00401"`
	InterchangeControlNumber          int64  `validate:"min=1,max=999999999"`
	AcknowledgementRequested          int    `validate:"oneof=0 1"`
	UsageIndicator                    string `validate:"oneof=P T"`
	ComponentElementSeparator         string `validate:"eq=0x7C"` // Have to escape pipe symbol
}

// StringArray converts ISA to an array of strings
func (s *ISA) StringArray() []string {
	return []string{
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
		fmt.Sprintf("%09d", s.InterchangeControlNumber),
		strconv.Itoa(s.AcknowledgementRequested),
		s.UsageIndicator,
		s.ComponentElementSeparator,
	}
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
	s.InterchangeControlNumber, err = strconv.ParseInt(elements[12], 10, 64)
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
