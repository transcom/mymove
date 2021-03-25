package edi824

import edisegment "github.com/transcom/mymove/pkg/edi/segment"

type transactionSet struct {
	ST   edisegment.ST    // transaction set header (bump up counter for "ST" and create new transactionSet)
	BGN  edisegment.BGN   // beginning statement
	OTIs []edisegment.OTI // original transaction identifications
	TEDs []edisegment.TED // technical error descriptions
	SE   edisegment.SE    // transaction set trailer
}

type functionalGroupEnvelope struct {
	GS              edisegment.GS // functional group header (bump up counter for "GS" and create new functionalGroupEnvelope)
	TransactionSets []transactionSet
	GE              edisegment.GE // functional group trailer
}

type interchangeControlEnvelope struct {
	ISA              edisegment.ISA // interchange control header
	FunctionalGroups []functionalGroupEnvelope
	IEA              edisegment.IEA // interchange control trailer
}

// EDI holds all the segments to parse an EDI 997
type EDI struct {
	InterchangeControlEnvelope interchangeControlEnvelope
}
