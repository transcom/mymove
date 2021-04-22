package edi824

import (
	"bufio"
	"fmt"
	"strings"

	edisegment "github.com/transcom/mymove/pkg/edi/segment"
)

type transactionSetCounter struct {
	otiCounter int
	tedCounter int
}

type functionalGroupCounter struct {
	tsCounter int
	ts        []transactionSetCounter
}

// 	ISA > FGs > TSs > OTI|TED
type counterData struct {
	fgCounter int
	fg        []functionalGroupCounter
}

// Parse takes in a string representation of a 824 EDI file and reads it into a 824 EDI struct
func (e *EDI) Parse(ediString string) error {
	var err error
	counter := counterData{}

	scanner := bufio.NewScanner(strings.NewReader(ediString))
	for scanner.Scan() {
		record := strings.Split(scanner.Text(), "*")

		if len(record) == 0 || len(strings.TrimSpace(record[0])) == 0 {
			continue
		}
		switch record[0] {
		case "ISA":
			err = e.InterchangeControlEnvelope.ISA.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("824 failed to parse %w", err)
			}
		case "GS":
			// functional group header
			// bump up counter fgCounter
			// create new functionalGroupEnvelope
			// inside functional group
			fg := functionalGroupEnvelope{}
			err = fg.GS.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("824 failed to parse %w", err)
			}
			e.InterchangeControlEnvelope.FunctionalGroups = append(e.InterchangeControlEnvelope.FunctionalGroups, fg)
			counter.fgCounter++
			counter.fg = append(counter.fg, functionalGroupCounter{})
		case "ST":
			// bump up counter for tsCounter
			// create new TransactionSet
			// inside functional group > transaction set
			fgIndex := counter.fgCounter - 1
			ts := TransactionSet{}
			err = ts.ST.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("824 failed to parse %w", err)
			}
			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets = append(e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets, ts)
			counter.fg[fgIndex].tsCounter++
			counter.fg[fgIndex].ts = append(counter.fg[fgIndex].ts, transactionSetCounter{})
		case "BGN":
			// beginning statement
			// inside functional group > transaction set
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			err = e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].BGN.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("824 failed to parse %w", err)
			}
		case "OTI":
			// bump up counter for otiCounter
			// inside functional group > transaction set
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			oti := edisegment.OTI{}
			err = oti.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("824 failed to parse %w", err)
			}
			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].OTIs = append(e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].OTIs, oti)
			counter.fg[fgIndex].ts[tsIndex].otiCounter++
		case "TED":
			// bump up counter for tedCounter
			// inside functional group > transaction set
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1

			ted := edisegment.TED{}
			err = ted.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("824 failed to parse %w", err)
			}
			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].TEDs = append(e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].TEDs, ted)
			counter.fg[fgIndex].ts[tsIndex].tedCounter++
		case "SE": // trailer to ST
			// transaction set trailer
			// inside functional group > transaction set
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			err = e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].SE.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("824 failed to parse %w", err)
			}
		case "GE": // trailer to GS
			// functional group trailer
			// inside functional group
			fgIndex := counter.fgCounter - 1
			err = e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].GE.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("824 failed to parse %w", err)
			}
		case "IEA": // trailer to ISA
			err = e.InterchangeControlEnvelope.IEA.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("824 failed to parse %w", err)
			}
		default:
			return fmt.Errorf("unexpected row for EDI 824, do not know how to parse: %s with %d parts", strings.Join(record, " "), len(record))
		}
	} // end of scanner loop

	return nil
}
