package ediinvoice

import (
	"bytes"

	"github.com/transcom/mymove/pkg/edi"
	edisegment "github.com/transcom/mymove/pkg/edi/segment"
)

// ICNSequenceName used to query Interchange Control Numbers from DB
const ICNSequenceName = "interchange_control_number"

// ICNRandomMin is the smallest allowed random-number based ICN (we use random ICN numbers in development)
const ICNRandomMin int64 = 100000000

// ICNRandomMax is the largest allowed random-number based ICN (we use random ICN numbers in development)
const ICNRandomMax int64 = 999999999

// Invoice858C holds all the segments that are generated
type Invoice858C struct {
	ISA          edisegment.ISA
	GS           edisegment.GS
	Header       []edisegment.Segment
	ServiceItems []edisegment.Segment
	GE           edisegment.GE
	IEA          edisegment.IEA
}

// Segments returns the invoice as an array of rows (string arrays),
// each containing a segment, to prepare it for writing
func (invoice Invoice858C) Segments() [][]string {
	records := [][]string{
		invoice.ISA.StringArray(),
		invoice.GS.StringArray(),
	}

	for _, line := range invoice.Header {
		records = append(records, line.StringArray())
	}
	for _, line := range invoice.ServiceItems {
		records = append(records, line.StringArray())
	}
	records = append(records, invoice.GE.StringArray())
	records = append(records, invoice.IEA.StringArray())
	return records
}

// EDIString returns the EDI representation of an 858C
func (invoice Invoice858C) EDIString() (string, error) {
	var b bytes.Buffer
	ediWriter := edi.NewWriter(&b)
	err := ediWriter.WriteAll(invoice.Segments())
	if err != nil {
		return "", err
	}
	return b.String(), err
}
