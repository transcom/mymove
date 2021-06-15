package pricing

import (
	"fmt"

	"github.com/tealeg/xlsx/v3"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

var parseDomesticMoveAccessorialPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 17 // 5a) Access. and Add. Prices
	const domAccessorialRowIndexStart int = 11
	const firstColumnIndexStart = 2
	const secondColumnIndexStart = 3
	const thirdColumnIndexStart = 4

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticMoveAccessorialPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	logger.Info("Parsing domestic move accessorial prices")
	var prices []models.StageDomesticMoveAccessorialPrice

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIdx := domAccessorialRowIndexStart; rowIdx < sheet.MaxRow; rowIdx++ {
		price := models.StageDomesticMoveAccessorialPrice{
			ServicesSchedule: getCell(sheet, rowIdx, firstColumnIndexStart),
			ServiceProvided:  getCell(sheet, rowIdx, secondColumnIndexStart),
			PricePerUnit:     getCell(sheet, rowIdx, thirdColumnIndexStart),
		}

		// All the rows are consecutive, if we get a blank we're done
		if price.ServicesSchedule == "" {
			break
		}

		if params.ShowOutput {
			logger.Info("", zap.Any("StageDomesticMoveAccessorialPrice", price))
		}
		prices = append(prices, price)
	}
	return prices, nil
}

var parseInternationalMoveAccessorialPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 17 // 5a) Access. and Add. Prices
	const intlAccessorialRowIndexStart int = 25
	const firstColumnIndexStart = 2
	const secondColumnIndexStart = 3
	const thirdColumnIndexStart = 4

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseInternationalMoveAccessorialPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	logger.Info("Parsing international move accessorial prices")
	var prices []models.StageInternationalMoveAccessorialPrice

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIdx := intlAccessorialRowIndexStart; rowIdx < sheet.MaxRow; rowIdx++ {
		price := models.StageInternationalMoveAccessorialPrice{
			Market:          getCell(sheet, rowIdx, firstColumnIndexStart),
			ServiceProvided: getCell(sheet, rowIdx, secondColumnIndexStart),
			PricePerUnit:    getCell(sheet, rowIdx, thirdColumnIndexStart),
		}

		// All the rows are consecutive, if we get a blank we're done
		if price.Market == "" {
			break
		}

		if params.ShowOutput {
			logger.Info("", zap.Any("StageInternationalMoveAccessorialPrice", price))
		}
		prices = append(prices, price)
	}
	return prices, nil
}

var parseDomesticInternationalAdditionalPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 17 // 5a) Access. and Add. Prices
	const additionalPricesRowIndexStart int = 39
	const firstColumnIndexStart = 2
	const secondColumnIndexStart = 3
	const thirdColumnIndexStart = 4

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticInternationalAdditionalPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	logger.Info("Parsing domestic/international additional prices")
	var prices []models.StageDomesticInternationalAdditionalPrice

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIdx := additionalPricesRowIndexStart; rowIdx < sheet.MaxRow; rowIdx++ {
		price := models.StageDomesticInternationalAdditionalPrice{
			Market:       getCell(sheet, rowIdx, firstColumnIndexStart),
			ShipmentType: getCell(sheet, rowIdx, secondColumnIndexStart),
			Factor:       getCell(sheet, rowIdx, thirdColumnIndexStart),
		}

		// All the rows are consecutive, if we get a blank we're done
		if price.Market == "" {
			break
		}

		if params.ShowOutput {
			logger.Info("", zap.Any("StageDomesticInternationalAdditionalPrice", price))
		}
		prices = append(prices, price)
	}
	return prices, nil
}

var verifyAccessAndAddPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	// XLSX Sheet consts
	const xlsxDataSheetNum = 17 // 5a) Access. and Add. Prices
	const domAccessorialRowIndexStart = 11
	const intlAccessorialRowIndexStart = 25
	const additionalPricesRowIndexStart = 39

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyAccessAndAddPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]

	err := helperCheckHeadersFor5a("Services Schedule", "Service Provided", "PricePerUnitofMeasure", sheet, domAccessorialRowIndexStart-2, domAccessorialRowIndexStart-1)
	if err != nil {
		return err
	}

	err = helperCheckHeadersFor5a("X", "EXAMPLE (per unit of measure)", "$X.XX", sheet, domAccessorialRowIndexStart-1, domAccessorialRowIndexStart)
	if err != nil {
		return err
	}

	err = helperCheckHeadersFor5a("Market", "Service Provided", "PricePerUnitofMeasure", sheet, intlAccessorialRowIndexStart-2, intlAccessorialRowIndexStart-1)
	if err != nil {
		return err
	}

	err = helperCheckHeadersFor5a("X", "EXAMPLE (per unit of measure)", "$X.XX", sheet, intlAccessorialRowIndexStart-1, intlAccessorialRowIndexStart)
	if err != nil {
		return err
	}

	err = helperCheckHeadersFor5a("Market", "Shipment Type", "Factor", sheet, additionalPricesRowIndexStart-2, additionalPricesRowIndexStart-1)
	if err != nil {
		return err
	}

	return helperCheckHeadersFor5a("CONUS / OCONUS", "EXAMPLE", "X.XX", sheet, additionalPricesRowIndexStart-1, additionalPricesRowIndexStart)
}

func helperCheckHeadersFor5a(firstHeader string, secondHeader string, thirdHeader string, sheet *xlsx.Sheet, dataRowsIdxBegin, dataRowsIdxEnd int) error {
	const firstColumnIndexStart = 2
	const secondColumnIndexStart = 3
	const thirdColumnIndexStart = 4

	for rowIdx := dataRowsIdxBegin; rowIdx < dataRowsIdxEnd; rowIdx++ {
		if header := getCell(sheet, rowIdx, firstColumnIndexStart); header != firstHeader {
			return fmt.Errorf("verifyAccessAndAddPrices expected to find header '%s', but received header '%s'", firstHeader, header)
		}
		if header := getCell(sheet, rowIdx, secondColumnIndexStart); header != secondHeader {
			return fmt.Errorf("verifyAccessAndAddPrices expected to find header '%s', but received header '%s'", secondHeader, header)
		}
		if header := removeWhiteSpace(getCell(sheet, rowIdx, thirdColumnIndexStart)); header != thirdHeader {
			return fmt.Errorf("verifyAccessAndAddPrices expected to find header '%s', but received header '%s'", thirdHeader, header)
		}
	}
	return nil
}
