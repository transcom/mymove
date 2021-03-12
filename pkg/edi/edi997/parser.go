package edi997

import (
	"bufio"
	"fmt"
	"strings"

	edisegment "github.com/transcom/mymove/pkg/edi/segment"
)

// picture of what the envelopes look like https://docs.oracle.com/cd/E19398-01/820-1275/agdaw/index.html

type dataSegment struct {
	//AK3      edisegment.AK3 // data segment note (bump up counter for "AK3", create new dataSegment)
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
	//AK9          edisegment.AK9 // functional group response trailer
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
	// b := bytes.NewBufferString(ediString)

	counter := counterData{}

	scanner := bufio.NewScanner(strings.NewReader(ediString))
	for scanner.Scan() {
		fmt.Println(scanner.Text())

		/*
			ediReader := edi.NewReader(strings.NewReader(scanner.Text()))
			row, err := ediReader.Read()
		*/

		/*
			    Rows are not all of the same length, so this seems to fail expecting the same length for all
			    rows. Switching to ediReader.Read() also seems to not be a viable option
				ediReader := edi.NewReader(strings.NewReader(ediString))
				ediRows, err := ediReader.ReadAll()
			    for _, row := range ediRows {
			    }
		*/
		record := strings.Split(scanner.Text(), "*")

		if len(record) == 0 {
			continue
		}
		switch record[0] {
		case "ISA":
			e.InterchangeControlEnvelope.ISA.Parse(record[1:])
		case "GS":
			// functional group header
			// bump up counter fgCounter
			// create new functionalGroupEnvelope
			// inside functional group
			fg := functionalGroupEnvelope{}
			fg.GS.Parse(record[1:])
			e.InterchangeControlEnvelope.FunctionalGroups = append(e.InterchangeControlEnvelope.FunctionalGroups, fg)
			counter.fgCounter++
			counter.fg = append(counter.fg, functionalGroupCounter{})
		case "ST":
			// bump up counter for tsCounter
			// create new transactionSet
			// inside functional group > transaction set
			fgIndex := counter.fgCounter - 1
			ts := transactionSet{}
			ts.ST.Parse(record[1:])
			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets = append(e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets, ts)
			counter.fg[fgIndex].tsCounter++
			counter.fg[fgIndex].ts = append(counter.fg[fgIndex].ts, transactionSetCounter{})
		case "AK1":
			// create new functionalGroupResponse
			// inside functional group > transaction set > functional group response
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.AK1.Parse(record[1:])
		case "AK2":
			// bump up counter for tsrCounter
			// create new transactionSetResponse
			// inside functional group > transaction set > functional group response > transaction set response
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			tsr := transactionSetResponse{}
			tsr.AK2.Parse(record[1:])
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
			// ds.AK3.Parse(record[1:])
			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.TransactionSetResponses[tsrIndex].dataSegments = append(e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.TransactionSetResponses[tsrIndex].dataSegments, ds)

			counter.fg[fgIndex].ts[tsIndex].fgr.tsr[tsrIndex].dsCounter++
		case "AK4": // trailer to AK3
			// inside functional group > transaction set > functional group response > transaction set response > data segment
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			tsrIndex := counter.fg[fgIndex].ts[tsIndex].fgr.tsrCounter - 1
			dsIndex := counter.fg[fgIndex].ts[tsIndex].fgr.tsr[tsrIndex].dsCounter - 1

			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.TransactionSetResponses[tsrIndex].dataSegments[dsIndex].AK4.Parse(record[1:])
		case "AK5": // trailer to AK2
			// transaction set response
			// inside functional group > transaction set > functional group response > transaction set response
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			tsrIndex := counter.fg[fgIndex].ts[tsIndex].fgr.tsrCounter - 1

			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].FunctionalGroupResponse.TransactionSetResponses[tsrIndex].AK5.Parse(record[1:])
		case "AK9": // trailer to AK1
			// functional group response trailer
			// inside functional group > transaction set > functional group response
			//fgIndex := counter.fgCounter - 1
			//tsIndex := counter.fg[fgIndex].tsCounter - 1
			//e.interchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].functionalGroupResponse.AK9.Parse(record[1:])
		case "SE": // trailer to ST
			// transaction set trailer
			// inside functional group > transaction set
			fgIndex := counter.fgCounter - 1
			tsIndex := counter.fg[fgIndex].tsCounter - 1
			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].TransactionSets[tsIndex].SE.Parse(record[1:])
		case "GE": // trailer to GS
			// functional group trailer
			// inside functional group
			fgIndex := counter.fgCounter - 1
			e.InterchangeControlEnvelope.FunctionalGroups[fgIndex].GE.Parse(record[1:])
		case "IEA": // trailer to ISA
			e.InterchangeControlEnvelope.IEA.Parse(record[1:])
		}
	} // end of scanner loop

	return nil
}
