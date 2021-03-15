package edi997

import edisegment "github.com/transcom/mymove/pkg/edi/segment"

// EDI 997 is a Functional Acknowledgement message that is sent to the MilMove system to simply acknowledge
// that the corresponding EDI 858 message was received.

// Picture of what the envelopes look like https://docs.oracle.com/cd/E19398-01/820-1275/agdaw/index.html

type dataSegment struct {
	AK3 edisegment.AK3 // data segment note (bump up counter for "AK3", create new dataSegment)
	AK4 edisegment.AK4 // data element note
}

type transactionSetResponse struct {
	AK2          edisegment.AK2 // transaction set response header (bump up counter for "AK2", create new transactionSetResponse)
	dataSegments []dataSegment  // data segments, loop ID AK3
	AK5          edisegment.AK5 // transaction set response trailer
}

type functionalGroupResponse struct {
	AK1                     edisegment.AK1           // functional group response header (create new functionalGroupResponse)
	TransactionSetResponses []transactionSetResponse // transaction set responses, loop ID AK2
	AK9                     edisegment.AK9           // functional group response trailer
}

type transactionSet struct {
	ST                      edisegment.ST // transaction set header (bump up counter for "ST" and create new transactionSet)
	FunctionalGroupResponse functionalGroupResponse
	SE                      edisegment.SE // transaction set trailer
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
