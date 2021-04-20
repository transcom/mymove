package edi824

import (
	"github.com/go-playground/validator/v10"

	edisegment "github.com/transcom/mymove/pkg/edi/segment"
)

// TransactionSet holds the transaction set envelope for the 824
type TransactionSet struct {
	ST   edisegment.ST    // transaction set header (bump up counter for "ST" and create new TransactionSet)
	BGN  edisegment.BGN   // beginning statement
	OTIs []edisegment.OTI `validate:"min=1,dive"` // original transaction identifications
	TEDs []edisegment.TED `validate:"dive"`       // technical error descriptions
	SE   edisegment.SE    // transaction set trailer
}

type functionalGroupEnvelope struct {
	GS              edisegment.GS    // functional group header (bump up counter for "GS" and create new functionalGroupEnvelope)
	TransactionSets []TransactionSet `validate:"min=1,dive"`
	GE              edisegment.GE    // functional group trailer
}

type interchangeControlEnvelope struct {
	ISA              edisegment.ISA            // interchange control header
	FunctionalGroups []functionalGroupEnvelope `validate:"min=1,dive"`
	IEA              edisegment.IEA            // interchange control trailer
}

// EDI holds all the segments to parse an EDI 824
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
func (edi824 EDI) Validate() error {
	return validate.Struct(edi824.InterchangeControlEnvelope)
}
