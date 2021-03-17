package edi997

import (
	"bufio"
	"fmt"
	"strings"
)

type transactionSetResponseCounter struct {
	dsCounter int
}

type functionalGroupResponseCounter struct {
	tsrCounter int
	tsr        []transactionSetResponseCounter
}

type transactionSetCounter struct {
	fgr functionalGroupResponseCounter
}

type functionalGroupCounter struct {
	tsCounter int
	ts        []transactionSetCounter
}

// 	ISA > FGs > TSs > FGR > TSRs > DSs
type counterData struct {
	fgCounter int
	fg        []functionalGroupCounter
}

// Parse takes in a string representation of a 997 EDI file and reads it into a 997 EDI struct
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
				return fmt.Errorf("997 failed to parse %w", err)
			}
		case "GS":
			// functional group header
			// bump up counter fgCounter
			// create new functionalGroupEnvelope
			// inside functional group
			fg := functionalGroupEnvelope{}
			err = fg.GS.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("997 failed to parse %w", err)
			}
			e.InterchangeControlEnvelope.FunctionalGroups = append(e.InterchangeControlEnvelope.FunctionalGroups, fg)
			counter.fgCounter++
			counter.fg = append(counter.fg, functionalGroupCounter{})
		case "ST":
			// bump up counter for tsCounter
			// create new transactionSet
			// inside functional group > transaction set
			fgIndex := counter.fgCounter - 1
			ts := transactionSet{}
			err = ts.ST.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("997 failed to parse %w", err)
			}
			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets = append(e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets, ts)
			counter.fg[fgIndex].tsCounter++
			counter.fg[fgIndex].ts = append(counter.fg[fgIndex].ts, transactionSetCounter{})
		case "AK1":
			// create new functionalGroupResponse
			// inside functional group > transaction set > functional group response
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			err = e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.AK1.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("997 failed to parse %w", err)
			}
		case "AK2":
			// bump up counter for tsrCounter
			// create new transactionSetResponse
			// inside functional group > transaction set > functional group response > transaction set response
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			tsr := transactionSetResponse{}
			err = tsr.AK2.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("997 failed to parse %w", err)
			}
			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.TransactionSetResponses = append(e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.TransactionSetResponses, tsr)
			counter.fg[fgIndex].ts[tsIndex].fgr.tsrCounter++
			counter.fg[fgIndex].ts[tsIndex].fgr.tsr = append(counter.fg[fgIndex].ts[tsIndex].fgr.tsr, transactionSetResponseCounter{})
		case "AK3":
			// bump up counter for dsCounter
			// create new dataSegment
			// inside functional group > transaction set > functional group response > transaction set response > data segment
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			tsrIndex := counter.fg[fgIndex].ts[tsIndex].fgr.tsrCounter - 1

			ds := dataSegment{}
			err = ds.AK3.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("997 failed to parse %w", err)
			}
			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.TransactionSetResponses[tsrIndex].dataSegments = append(e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.TransactionSetResponses[tsrIndex].dataSegments, ds)

			counter.fg[fgIndex].ts[tsIndex].fgr.tsr[tsrIndex].dsCounter++
		case "AK4": // trailer to AK3
			// inside functional group > transaction set > functional group response > transaction set response > data segment
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			tsrIndex := counter.fg[fgIndex].ts[tsIndex].fgr.tsrCounter - 1
			dsIndex := counter.fg[fgIndex].ts[tsIndex].fgr.tsr[tsrIndex].dsCounter - 1

			err = e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.TransactionSetResponses[tsrIndex].dataSegments[dsIndex].AK4.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("997 failed to parse %w", err)
			}
		case "AK5": // trailer to AK2
			// transaction set response
			// inside functional group > transaction set > functional group response > transaction set response
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			tsrIndex := counter.fg[fgIndex].ts[tsIndex].fgr.tsrCounter - 1

			err = e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.TransactionSetResponses[tsrIndex].AK5.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("997 failed to parse %w", err)
			}
		case "AK9": // trailer to AK1
			// functional group response trailer
			// inside functional group > transaction set > functional group response
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			err = e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.AK9.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("997 failed to parse %w", err)
			}
		case "SE": // trailer to ST
			// transaction set trailer
			// inside functional group > transaction set
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			err = e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].SE.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("997 failed to parse %w", err)
			}
		case "GE": // trailer to GS
			// functional group trailer
			// inside functional group
			fgIndex := counter.fgCounter - 1
			err = e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].GE.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("997 failed to parse %w", err)
			}
		case "IEA": // trailer to ISA
			err = e.InterchangeControlEnvelope.IEA.Parse(record[1:])
			if err != nil {
				return fmt.Errorf("997 failed to parse %w", err)
			}
		default:
			return fmt.Errorf("unexpected row for EDI 997, do not know how to parse: %s with %d parts", strings.Join(record, " "), len(record))
		}
	} // end of scanner loop

	return nil
}
