package pricing

import (
	"fmt"
	"log"

	"github.com/transcom/mymove/pkg/models"
)

// parseOtherIntlPrices: parser for: 3d) Other International Prices
var parseOtherIntlPrices processXlsxSheet = func(params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 13 // 3d) International Other Prices
	const feeColIndexStart int = 4  // start at column 6 to get the rates
	const feeRowIndexStart int = 10 // start at row 10 to get the rates
	const priceAreaCodeColumn int = 2
	const priceAreaNameColumn int = 3

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseOtherIntlPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	var otherIntlPrices []models.StageOtherIntlPrice
	dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowIndexStart:]

	for _, row := range dataRows {
		colIndex := feeColIndexStart
		// All the rows are consecutive, if we get to a blank one we're done
		if getCell(row.Cells, colIndex) == "" {
			break
		}

		for _, s := range rateSeasons {
			otherIntlPrice := models.StageOtherIntlPrice{
				RateAreaCode: getCell(row.Cells, priceAreaCodeColumn),
				RateAreaName: getCell(row.Cells, priceAreaNameColumn),
				Season:       s,
			}
			otherIntlPrice.HHGOriginPackPrice = getCell(row.Cells, colIndex)
			colIndex++
			otherIntlPrice.HHGDestinationUnPackPrice = getCell(row.Cells, colIndex)
			colIndex++
			otherIntlPrice.UBOriginPackPrice = getCell(row.Cells, colIndex)
			colIndex++
			otherIntlPrice.UBDestinationUnPackPrice = getCell(row.Cells, colIndex)
			colIndex++
			otherIntlPrice.OriginDestinationSITFirstDayWarehouse = getCell(row.Cells, colIndex)
			colIndex++
			otherIntlPrice.OriginDestinationSITAddlDays = getCell(row.Cells, colIndex)
			colIndex++
			otherIntlPrice.SITLte50Miles = getCell(row.Cells, colIndex)
			colIndex++
			otherIntlPrice.SITGt50Miles = getCell(row.Cells, colIndex)
			colIndex += 2

			if params.ShowOutput == true {
				log.Printf("%v\n", otherIntlPrice)

			}
			otherIntlPrices = append(otherIntlPrices, otherIntlPrice)
		}
	}

	return otherIntlPrices, nil
}

var verifyOtherIntlPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 13 // 3d) International Other Prices
	const feeColIndexStart int = 4  // start at column 6 to get the rates
	const feeRowIndexStart int = 10 // start at row 10 to get the rates
	const priceAreaCodeColumn int = 2
	const priceAreaNameColumn int = 3

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyOtherIntlPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	return nil
}
