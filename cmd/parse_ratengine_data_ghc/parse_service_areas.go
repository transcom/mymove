package main

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop"
)

// parseServiceAreas: parser for: 1b) Service Areas
var parseServiceAreas processXlsxSheet = func(params paramConfig, sheetIndex int, db *pop.Connection) error {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 4          // 1b) Service Areas
	const serviceAreaRowIndexStart int = 10 // start at row 10 to get the rates
	const basePointCityColumn int = 2
	const stateColumn int = 3
	const serviceAreaNumberColumn int = 4
	const zip3sColumn int = 5
	const internationalRateAreaColumn int = 9
	const rateAreaIDColumn int = 10

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("parseServiceAreas expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	log.Println("Parsing Domestic Service Areas")
	// Create CSV writer to save data to CSV file, returns nil if params.saveToFile=false
	csvWriter := createCsvWriter(params.saveToFile, sheetIndex, params.runTime)
	if csvWriter != nil {
		defer csvWriter.close()

		// Write header to CSV
		dsa := domesticServiceArea{}
		csvWriter.write(dsa.csvHeader())
	}

	dataRows := params.xlsxFile.Sheets[xlsxDataSheetNum].Rows[serviceAreaRowIndexStart:]
	for _, row := range dataRows {
		domServArea := domesticServiceArea{
			BasePointCity:     getCell(row.Cells, basePointCityColumn),
			State:             getCell(row.Cells, stateColumn),
			ServiceAreaNumber: formatServiceAreaNumber(getCell(row.Cells, serviceAreaNumberColumn)),
			Zip3s:             splitZip3s(getCell(row.Cells, zip3sColumn)),
		}
		// All the rows are consecutive, if we get to a blank one we're done
		if domServArea.BasePointCity == "" {
			break
		} else if csvWriter != nil {
			csvWriter.write(domServArea.toSlice())
		}
		domServArea.saveToDatabase(db)
	}

	log.Println("Parsing International Service Areas")
	// Create CSV writer to save data to CSV file, returns nil if params.saveToFile=false
	if csvWriter != nil {
		// Write header to CSV
		isa := internationalServiceArea{}
		csvWriter.write(isa.csvHeader())
	}

	for _, row := range dataRows {
		intlServArea := internationalServiceArea{
			RateArea:   getCell(row.Cells, internationalRateAreaColumn),
			RateAreaID: getCell(row.Cells, rateAreaIDColumn),
		}
		// All the rows are consecutive, if we get to a blank one we're done
		if intlServArea.RateArea == "" {
			break
		} else if csvWriter != nil {
			csvWriter.write(intlServArea.toSlice())
		}
	}
	return nil
}

// verifyServiceAreas: verification for: 1b) Service Areas
var verifyServiceAreas verifyXlsxSheet = func(params paramConfig, sheetIndex int) error {
	log.Println("TODO verifyServiceAreas() not implemented")
	return nil
}
