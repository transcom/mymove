package edi

import (
	"fmt"
	"io"

	"github.com/transcom/mymove/pkg/edi/invoice"
	edi_segment "github.com/transcom/mymove/pkg/edi/segment"
)

// TODO: Get these from the first line of the file?
const fieldSep = '*'

// ParseEDIFile takes a file as input and prints the parsed file (for now)
func ParseEDIFile(r io.Reader) error {
	scanner := NewScanner(r, fieldSep)
	var transaction []edi_segment.Segment
	for scanner.Next() {
		segment := scanner.Segment()
		switch segment.(type) {
		case *edi_segment.ISA:
		case *edi_segment.GS:
		case *edi_segment.ST:
			transaction = []edi_segment.Segment{segment}
		case *edi_segment.SE:
			// TODO: determine which parser to use based on GS segment
			parser := invoice.NewParser859(append(transaction, segment))
			err := parser.Parse()
			if err != nil {
				return err
			}
			fmt.Printf("%#v\n", parser.Invoice())
		case *edi_segment.GE:
		case *edi_segment.IEA:
		default:
			transaction = append(transaction, segment)
		}
	}

	return scanner.Err()
}
