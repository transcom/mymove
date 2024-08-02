package tppspaidinvoicereport

import edisegment "github.com/transcom/mymove/pkg/edi/segment"

// look to pkg/edi/edi824/edi.go for reference

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

// EDI holds all the segments to parse TPPS paid invoice report
type EDI struct {
	InterchangeControlEnvelope interchangeControlEnvelope
}
