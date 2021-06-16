package pricing

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

// parseOtherIntlPrices: parser for: 3d) Other International Prices
var parseOtherIntlPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 13 // 3d) International Other Prices
	const feeColIndexStart int = 4  // start at column 6 to get the rates
	const feeRowIndexStart int = 10 // start at row 10 to get the rates
	const priceAreaCodeColumn int = 2
	const priceAreaNameColumn int = 3

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseOtherIntlPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	logger.Info("Parsing other international prices")

	var otherIntlPrices []models.StageOtherIntlPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := feeRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		colIndex := feeColIndexStart
		// All the rows are consecutive, if we get to a blank one we're done
		if getCell(sheet, rowIndex, colIndex) == "" {
			break
		}

		for _, s := range rateSeasons {
			otherIntlPrice := models.StageOtherIntlPrice{
				RateAreaCode: getCell(sheet, rowIndex, priceAreaCodeColumn),
				RateAreaName: getCell(sheet, rowIndex, priceAreaNameColumn),
				Season:       s,
			}
			otherIntlPrice.HHGOriginPackPrice = getCell(sheet, rowIndex, colIndex)
			colIndex++
			otherIntlPrice.HHGDestinationUnPackPrice = getCell(sheet, rowIndex, colIndex)
			colIndex++
			otherIntlPrice.UBOriginPackPrice = getCell(sheet, rowIndex, colIndex)
			colIndex++
			otherIntlPrice.UBDestinationUnPackPrice = getCell(sheet, rowIndex, colIndex)
			colIndex++
			otherIntlPrice.OriginDestinationSITFirstDayWarehouse = getCell(sheet, rowIndex, colIndex)
			colIndex++
			otherIntlPrice.OriginDestinationSITAddlDays = getCell(sheet, rowIndex, colIndex)
			colIndex++
			otherIntlPrice.SITLte50Miles = getCell(sheet, rowIndex, colIndex)
			colIndex++
			otherIntlPrice.SITGt50Miles = getCell(sheet, rowIndex, colIndex)
			colIndex += 2

			if params.ShowOutput {
				logger.Info("", zap.Any("StageOtherIntlPrice", otherIntlPrice))
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

	repeatingHeaders := []string{
		"HHGOriginPackPrice(percwt)",
		"HHGDestinationUnpackPrice(percwt)",
		"UBOriginPackPrice(percwt)",
		"UBDestinationUnpackPrice(percwt)",
		"Origin/DestinationSITFirstDay&WarehouseHandling(percwt)",
		"Origin/DestinationSITAdd'lDays(percwt)",
		"SITPickup/Delivery≤50Miles(percwt)",
		"SITPickup/Delivery>50Miles(percwtpermile)",
	}

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyOtherIntlPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]

	nonPriceHeaderIndex := feeRowIndexStart - 3
	headerRowIndex := feeRowIndexStart - 2

	if err := verifyHeader(sheet, nonPriceHeaderIndex, priceAreaCodeColumn, "PriceAreaCode/ID"); err != nil {
		return fmt.Errorf("verifyOtherIntlPrices verification failure: %w", err)

	}

	priceAreaNameHeader := "InternationalPriceArea(PPIRA)/DomesticPriceArea(PPDRA)/Non-StandardRateArea"
	if err := verifyHeader(sheet, nonPriceHeaderIndex, priceAreaNameColumn, priceAreaNameHeader); err != nil {
		return fmt.Errorf("verifyOtherIntlPrices verification failure: %w", err)
	}

	// NonPeak season headers
	colIndex := feeColIndexStart
	for _, repeatingHeader := range repeatingHeaders {
		if err := verifyHeader(sheet, headerRowIndex, colIndex, repeatingHeader); err != nil {
			return fmt.Errorf("verifyOtherIntlPrices verification failure: %w", err)
		}
		colIndex++
	}
	colIndex++

	// Peak season headers
	for _, repeatingHeader := range repeatingHeaders {
		if err := verifyHeader(sheet, headerRowIndex, colIndex, repeatingHeader); err != nil {
			return fmt.Errorf("verifyOtherIntlPrices verification failure: %w", err)
		}
		colIndex++
	}

	return nil
}
