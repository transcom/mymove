package edi997

import (
	"github.com/go-playground/validator/v10"

	edisegment "github.com/transcom/mymove/pkg/edi/segment"
)

// EDI 997 is a Functional Acknowledgement message that is sent to the MilMove system to simply acknowledge
// that the corresponding EDI 858 message was received.

// Picture of what the envelopes look like https://docs.oracle.com/cd/E19398-01/820-1275/agdaw/index.html

type dataSegment struct {
	AK3 edisegment.AK3 // data segment note (bump up counter for "AK3", create new dataSegment)
	AK4 edisegment.AK4 // data element note
}

type transactionSetResponse struct {
	AK2          edisegment.AK2 // transaction set response header (bump up counter for "AK2", create new transactionSetResponse)
	dataSegments []dataSegment  `validate:"dive"` // data segments, loop ID AK3
	AK5          edisegment.AK5 // transaction set response trailer
}

type functionalGroupResponse struct {
	AK1                     edisegment.AK1           // functional group response header (create new functionalGroupResponse)
	TransactionSetResponses []transactionSetResponse `validate:"dive"` // transaction set responses, loop ID AK2
	AK9                     edisegment.AK9           // functional group response trailer
}

type transactionSet struct {
	ST                      edisegment.ST // transaction set header (bump up counter for "ST" and create new transactionSet)
	FunctionalGroupResponse functionalGroupResponse
	SE                      edisegment.SE // transaction set trailer
}

type functionalGroupEnvelope struct {
	GS              edisegment.GS    // functional group header (bump up counter for "GS" and create new functionalGroupEnvelope)
	TransactionSets []transactionSet `validate:"min=1,dive"`
	GE              edisegment.GE    // functional group trailer
}

type interchangeControlEnvelope struct {
	ISA              edisegment.ISA            // interchange control header
	FunctionalGroups []functionalGroupEnvelope `validate:"min=1,dive"`
	IEA              edisegment.IEA            // interchange control trailer
}

// EDI holds all the segments to parse an EDI 997
type EDI struct {
	InterchangeControlEnvelope interchangeControlEnvelope
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate will validate the EDI997 (and nested structs) to make sure they will produce legal EDI.
// This returns either an InvalidValidationError or a validator.ValidationErrors that allows all validation
// errors to be introspected individually.
func (edi997 EDI) Validate() error {
	return validate.Struct(edi997.InterchangeControlEnvelope)
}
