package main

import (
	"fmt"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// parseDomesticServiceAreaPrices: parser for: 2b) Dom. Service Area Prices
var parseIntlOtherPrices processXlsxSheet = func(params paramConfig, sheetIndex int, tableFromSliceCreator services.TableFromSliceCreator, csvWriter *createCsvHelper) error {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 13  // 3d) International Other Prices
	const feeColIndexStart int = 4  // start at column 6 to get the rates
	const feeRowIndexStart int = 10 // start at row 10 to get the rates
	const priceAreaCodeColumn int = 2
	const priceAreaNameColumn int = 3

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("parseIntlOtherPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	var intlOtherPrices []models.StageIntlOtherPrice
	dataRows := params.xlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowIndexStart:]

	for _, row := range dataRows {
		colIndex := feeColIndexStart
		for _, s := range rateSeasons {
			intlOtherPrice := models.StageIntlOtherPrice {
				RateAreaCode: getCell(row.Cells, priceAreaCodeColumn),
				RateAreaName: getCell(row.Cells, priceAreaNameColumn),
				Season: s,
			}
			intlOtherPrice.HHGOriginPackPrice = getCell(row.Cells, colIndex)
			colIndex++
			intlOtherPrice.HHGDestinationUnPackPrice = getCell(row.Cells, colIndex)
			colIndex++
			intlOtherPrice.UBOriginPackPrice = getCell(row.Cells, colIndex)
			colIndex++
			intlOtherPrice.UBDestinationUnPackPrice = getCell(row.Cells, colIndex)
			colIndex++
			intlOtherPrice.OriginDestinationSITFirstDayWarehouse = getCell(row.Cells, colIndex)
			colIndex++
			intlOtherPrice.OriginDestinationSITAddlDays = getCell(row.Cells, colIndex)
			colIndex++
			intlOtherPrice.SITLte50Miles = getCell(row.Cells, colIndex)
			colIndex++
			intlOtherPrice.SITGt50Miles = getCell(row.Cells, colIndex)
			colIndex += 2

			intlOtherPrices = append(intlOtherPrices, intlOtherPrice)
		}

	}
}
