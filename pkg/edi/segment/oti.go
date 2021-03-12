package edisegment

import (
	"fmt"
	"strconv"
)

// OTI represents the OTI EDI segment
type OTI struct {
	ApplicationAcknowledgementCode   string `validate:"oneof=TA TE TR"`
	ReferenceIdentificationQualifier string `validate:"oneof=BM CN"`
	ReferenceIdentification          string `validate:"min=1,max=30"`
	ApplicationSendersCode           string `validate:"omitempty,min=2,max=15"`
	ApplicationReceiversCode         string `validate:"omitempty,min=2,max=15"`
	Date                             string `validate:"omitempty,datetime=20060102"`
	Time                             string `validate:"omitempty,datetime=1504"`
	GroupControlNumber               int64  `validate:"required_with=TransactionSetControlNumber,omitempty,min=1,max=999999999"`
	TransactionSetControlNumber      string `validate:"omitempty,min=4,max=9"`
}

// StringArray converts OTI to an array of strings
func (s *OTI) StringArray() []string {
	// For the optional int fields, make sure we map a zero to the empty string.
	var strGroupControlNumber string
	if s.GroupControlNumber != 0 {
		strGroupControlNumber = strconv.FormatInt(s.GroupControlNumber, 10)
	}

	return []string{
		"OTI",
		s.ApplicationAcknowledgementCode,
		s.ReferenceIdentificationQualifier,
		s.ReferenceIdentification,
		s.ApplicationSendersCode,
		s.ApplicationReceiversCode,
		s.Date,
		s.Time,
		strGroupControlNumber,
		s.TransactionSetControlNumber,
	}
}

// Parse parses an X12 string that's split into an array into the OTI struct
func (s *OTI) Parse(elements []string) error {
	expectedNumElements := 9
	if len(elements) != expectedNumElements {
		return fmt.Errorf("OTI: Wrong number of fields, expected %d, got %d", expectedNumElements, len(elements))
	}

	s.ApplicationAcknowledgementCode = elements[0]
	s.ReferenceIdentificationQualifier = elements[1]
	s.ReferenceIdentification = elements[2]
	s.ApplicationSendersCode = elements[3]
	s.ApplicationReceiversCode = elements[4]
	s.Date = elements[5]
	s.Time = elements[6]

	// For the optional int fields, make sure we map an empty string to a zero.
	var err error
	strGroupControlNumber := elements[7]
	if strGroupControlNumber == "" {
		s.GroupControlNumber = 0
	} else {
		s.GroupControlNumber, err = strconv.ParseInt(strGroupControlNumber, 10, 64)
		if err != nil {
			return err
		}
	}

	s.TransactionSetControlNumber = elements[8]

	return nil
}
