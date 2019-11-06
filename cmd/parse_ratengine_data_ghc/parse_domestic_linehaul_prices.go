package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// parseDomesticLinehaulPrices: parser for 2a) Domestic Linehaul Prices
var parseDomesticLinehaulPrices processXlsxSheet = func(params paramConfig, sheetIndex int, tableFromSliceCreator services.TableFromSliceCreator, csvWriter *createCsvHelper) error {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 6  // 2a) Domestic Linehaul Prices
	const feeColIndexStart int = 6  // start at column 6 to get the rates
	const feeRowIndexStart int = 14 // start at row 14 to get the rates
	const serviceAreaNumberColumn int = 2
	const originServiceAreaColumn int = 3
	const serviceScheduleColumn int = 4
	const numEscalationYearsToProcess = sharedNumEscalationYearsToProcess

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("parseDomesticLinehaulPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	var domPrices []models.StageDomesticLinehaulPrice
	dataRows := params.xlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowIndexStart:]
	for _, row := range dataRows {
		colIndex := feeColIndexStart
		// For number of baseline + Escalation years
		for escalation := 0; escalation < numEscalationYearsToProcess; escalation++ {
			// For each Rate Season
			for _, r := range rateSeasons {
				// For each weight band
				for _, w := range dLhWeightBands {
					// For each mileage range
					for _, m := range dLhMilesRanges {
						domPrice := models.StageDomesticLinehaulPrice{
							ServiceAreaNumber: getCell(row.Cells, serviceAreaNumberColumn),
							OriginServiceArea: getCell(row.Cells, originServiceAreaColumn),
							ServicesSchedule:  getCell(row.Cells, serviceScheduleColumn),
							Season:            r,
							WeightLower:       strconv.Itoa(w.lowerLbs),
							WeightUpper:       strconv.Itoa(w.upperLbs),
							MilesLower:        strconv.Itoa(m.lower),
							MilesUpper:        strconv.Itoa(m.upper),
							EscalationNumber:  strconv.Itoa(escalation),
							Rate:              getCell(row.Cells, colIndex),
						}
						colIndex++
						if params.showOutput == true {
							log.Printf("%v\n", domPrice)
						}
						domPrices = append(domPrices, domPrice)
					}
				}
				colIndex++ // skip 1 column (empty column) before starting next Rate type
			}
		}
	}

	// TODO: Move these two things out of here
	if csvWriter != nil {
		if err := gocsv.MarshalFile(domPrices, csvWriter.csvFile); err != nil {
			return errors.Wrap(err, "Could not marshal CSV file for domestic linehaul prices")
		}
	}
	if err := tableFromSliceCreator.CreateTableFromSlice(domPrices); err != nil {
		return errors.Wrap(err, "Could not create temp table for domestic linehaul prices")
	}

	return nil
}

// verifyDomesticLinehaulPrices: verification for 2a) Domestic Linehaul Prices
var verifyDomesticLinehaulPrices verifyXlsxSheet = func(params paramConfig, sheetIndex int) error {

	if dLhWeightBandNumCells != dLhWeightBandNumCellsExpected {
		return fmt.Errorf("parseDomesticLinehaulPrices(): Exepected %d columns per weight band, found %d defined in golang parser", dLhWeightBandNumCellsExpected, dLhWeightBandNumCells)
	}

	if len(dLhWeightBands) != dLhWeightBandCountExpected {
		return fmt.Errorf("parseDomesticLinehaulPrices(): Exepected %d weight bands, found %d defined in golang parser", dLhWeightBandCountExpected, len(dLhWeightBands))
	}

	// XLSX Sheet consts
	const xlsxDataSheetNum int = 6  // 2a) Domestic Linehaul Prices
	const feeColIndexStart int = 6  // start at column 6 to get the rates
	const feeRowIndexStart int = 14 // start at row 14 to get the rates
	const serviceAreaNumberColumn int = 2
	const originServiceAreaColumn int = 3
	const serviceScheduleColumn int = 4
	const numEscalationYearsToProcess int = 2

	// Check headers
	const feeRowMilageHeaderIndexStart = feeRowIndexStart - 3
	const verifyHeaderIndexEnd = feeRowMilageHeaderIndexStart + 2

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyDomesticLinehaulPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	dataRows := params.xlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowMilageHeaderIndexStart:verifyHeaderIndexEnd]
	for dataRowsIndex, row := range dataRows {
		colIndex := feeColIndexStart
		// For number of baseline + Escalation years
		for escalation := 0; escalation < numEscalationYearsToProcess; escalation++ {
			// For each Rate Season
			for _, r := range rateSeasons {
				// For each weight band
				for _, w := range dLhWeightBands {
					// For each milage range
					for dLhMilesRangesIndex, m := range dLhMilesRanges {
						// skip the last index because the text is not easily checked
						if dLhMilesRangesIndex == len(dLhMilesRanges)-1 {
							colIndex++
							continue
						}
						verificationLog := fmt.Sprintf(" , verfication for row index: %d, colIndex: %d, Escalation: %d, rateSeasons %v, dLhWeightBands %v",
							dataRowsIndex, colIndex, escalation, r, w)
						if dataRowsIndex == 0 {
							if m.lower != getInt(getCell(row.Cells, colIndex)) {
								return fmt.Errorf("format error: From Miles --> does not match expected number expected %d got %s\n%s", m.lower, getCell(row.Cells, colIndex), verificationLog)
							}
							if "ServiceAreaNumber" != removeWhiteSpace(getCell(row.Cells, serviceAreaNumberColumn)) {
								return fmt.Errorf("format error: Header <ServiceAreaNumber> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(row.Cells, serviceAreaNumberColumn)), verificationLog)
							}
							if "OriginServiceArea" != removeWhiteSpace(getCell(row.Cells, originServiceAreaColumn)) {
								return fmt.Errorf("format error: Header <OriginServiceArea> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(row.Cells, originServiceAreaColumn)), verificationLog)
							}
							if "ServicesSchedule" != removeWhiteSpace(getCell(row.Cells, serviceScheduleColumn)) {
								return fmt.Errorf("format error: Header <SServicesSchedule> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(row.Cells, serviceScheduleColumn)), verificationLog)
							}
						} else if dataRowsIndex == 1 {
							if m.upper != getInt(getCell(row.Cells, colIndex)) {
								return fmt.Errorf("format error: To Miles --> does not match expected number expected %d got %s\n%s", m.upper, getCell(row.Cells, colIndex), verificationLog)
							}
						} else if dataRowsIndex == 2 {
							if "EXAMPLE" != getCell(row.Cells, originServiceAreaColumn) {
								return fmt.Errorf("format error: Filler text <EXAMPLE> is missing got <%s> instead\n%s", getCell(row.Cells, originServiceAreaColumn), verificationLog)
							}
						}
						colIndex++
					}
				}
				colIndex++ // skip 1 column (empty column) before starting next Rate type
			}
		}
	}

	return nil
}
