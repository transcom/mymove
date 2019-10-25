package main

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop"
)

// parseDomesticServiceAreaPrices: parser for: 2b) Dom. Service Area Prices
var parseDomesticServiceAreaPrices processXlsxSheet = func(params paramConfig, sheetIndex int, db *pop.Connection) error {
	// Create CSV writer to save data to CSV file, returns nil if params.saveToFile=false
	csvWriter := createCsvWriter(params.saveToFile, sheetIndex, params.runTime)
	if csvWriter != nil {
		defer csvWriter.close()

		// Write header to CSV
		dp := domesticServiceAreaPrice{}
		csvWriter.write(dp.csvHeader())
	}

	// XLSX Sheet consts
	const xlsxDataSheetNum int = 7  // 2b) Domestic Service Area Prices
	const feeColIndexStart int = 6  // start at column 6 to get the rates
	const feeRowIndexStart int = 10 // start at row 10 to get the rates
	const serviceAreaNumberColumn int = 2
	const serviceAreaNameColumn int = 3
	const serviceScheduleColumn int = 4
	const sITPickupDeliveryScheduleColumn int = 5
	const numEscalationYearsToProcess = sharedNumEscalationYearsToProcess

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("parseDomesticServiceAreaPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	dataRows := params.xlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowIndexStart:]
	for _, row := range dataRows {
		colIndex := feeColIndexStart
		// For number of baseline + Escalation years
		for escalation := 0; escalation < numEscalationYearsToProcess; escalation++ {
			// For each Rate Season
			for _, r := range rateSeasons {
				domPrice := domesticServiceAreaPrice{
					ServiceAreaNumber:         formatServiceAreaNumber(getCell(row.Cells, serviceAreaNumberColumn)),
					ServiceAreaName:           getCell(row.Cells, serviceAreaNameColumn),
					ServiceSchedule:           getInt(getCell(row.Cells, serviceScheduleColumn)),
					SITPickupDeliverySchedule: getInt(getCell(row.Cells, sITPickupDeliveryScheduleColumn)),
					Season:                    r,
					Escalation:                escalation,
				}

				domPrice.ShorthaulPrice = removeFirstDollarSign(getCell(row.Cells, colIndex))
				colIndex++
				domPrice.OriginDestinationPrice = removeFirstDollarSign(getCell(row.Cells, colIndex))
				colIndex += 3 // skip 2 columns pack and unpack
				domPrice.OriginDestinationSITFirstDayWarehouse = removeFirstDollarSign(getCell(row.Cells, colIndex))
				colIndex++
				domPrice.OriginDestinationSITAddlDays = removeFirstDollarSign(getCell(row.Cells, colIndex))
				colIndex++ // skip column SIT Pickup / Delivery ≤50 miles (per cwt)

				if params.showOutput == true {
					log.Println(domPrice.toSlice())
				}
				if csvWriter != nil {
					csvWriter.write(domPrice.toSlice())
				}

				colIndex += 2 // skip 1 column (empty column) before starting next Rate type
			}

		}
	}

	return nil
}

// verifyDomesticServiceAreaPrices: verification 2b) Dom. Service Area Prices
var verifyDomesticServiceAreaPrices verifyXlsxSheet = func(params paramConfig, sheetIndex int) error {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 7  // 2a) Domestic Linehaul Prices
	const feeColIndexStart int = 6  // start at column 6 to get the rates
	const feeRowIndexStart int = 10 // start at row 10 to get the rates
	const serviceAreaNumberColumn int = 2
	const serviceAreaNameColumn int = 3
	const serviceScheduleColumn int = 4
	const sITPickupDeliveryScheduleColumn int = 5
	const numEscalationYearsToProcess int = 4

	// Check headers
	const feeRowMilageHeaderIndexStart = feeRowIndexStart - 2
	const verifyHeaderIndexEnd = feeRowMilageHeaderIndexStart + 2

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyDomesticServiceAreaPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	// Verify header strings
	repeatingHeaders := []string{
		"Shorthaul Price (per cwt per mile)",
		"Origin / Destination Price (per cwt)",
		"Origin Pack Price (per cwt)",
		"Destination Unpack Price (per cwt)",
		"Origin / Destination SIT First Day & Warehouse Handling (per cwt)",
		"Origin / Destination SIT Add'l Days (per cwt)",
		"SIT Pickup / Delivery ≤50 miles (per cwt)",
	}

	dataRows := params.xlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowMilageHeaderIndexStart:verifyHeaderIndexEnd]
	for dataRowsIndex, row := range dataRows {
		colIndex := feeColIndexStart
		// For number of baseline + Escalation years
		for escalation := 0; escalation < numEscalationYearsToProcess; escalation++ {
			// For each Rate Season
			for _, r := range rateSeasons {
				verificationLog := fmt.Sprintf(" , verfication for row index: %d, colIndex: %d, Escalation: %d, rateSeasons %v",
					dataRowsIndex, colIndex, escalation, r)

				if dataRowsIndex == 0 {
					if "ServiceAreaNumber" != removeWhiteSpace(getCell(row.Cells, serviceAreaNumberColumn)) {
						return fmt.Errorf("format error: Header <ServiceAreaNumber> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(row.Cells, serviceAreaNumberColumn)), verificationLog)
					}
					if "ServiceAreaName" != removeWhiteSpace(getCell(row.Cells, serviceAreaNameColumn)) {
						return fmt.Errorf("format error: Header <ServiceAreaName> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(row.Cells, serviceAreaNameColumn)), verificationLog)
					}
					if "ServicesSchedule" != removeWhiteSpace(getCell(row.Cells, serviceScheduleColumn)) {
						return fmt.Errorf("format error: Header <ServicesSchedule> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(row.Cells, serviceScheduleColumn)), verificationLog)
					}

					if "SITPickup/DeliverySchedule" != removeWhiteSpace(getCell(row.Cells, sITPickupDeliveryScheduleColumn)) {
						return fmt.Errorf("format error: Header <SIT Pickup / Delivery> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(row.Cells, sITPickupDeliveryScheduleColumn)), verificationLog)
					}

					for _, repeatingHeader := range repeatingHeaders {
						if removeWhiteSpace(repeatingHeader) != removeWhiteSpace(getCell(row.Cells, colIndex)) {
							return fmt.Errorf("format error: Header contains <%s> is missing got <%s> instead\n%s", removeWhiteSpace(repeatingHeader), removeWhiteSpace(getCell(row.Cells, colIndex)), verificationLog)
						}
						colIndex++
					}
					colIndex++ // skip 1 column (empty column) before starting next Rate type
				} else if dataRowsIndex == 1 {
					if "EXAMPLE" != removeWhiteSpace(getCell(row.Cells, serviceAreaNameColumn)) {
						return fmt.Errorf("format error: Filler text <EXAMPLE> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(row.Cells, serviceAreaNameColumn)), verificationLog)
					}
				}

			}

		}
	}
	return nil
}
