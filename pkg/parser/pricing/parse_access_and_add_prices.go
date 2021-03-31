package pricing

import (
	"fmt"

	"github.com/tealeg/xlsx"
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
	dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[domAccessorialRowIndexStart:]
	for _, row := range dataRows {
		price := models.StageDomesticMoveAccessorialPrice{
			ServicesSchedule: getCell(row.Cells, firstColumnIndexStart),
			ServiceProvided:  getCell(row.Cells, secondColumnIndexStart),
			PricePerUnit:     getCell(row.Cells, thirdColumnIndexStart),
		}

		// All the rows are consecutive, if we get a blank we're done
		if price.ServicesSchedule == "" {
			break
		}

		if params.ShowOutput == true {
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
	dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[intlAccessorialRowIndexStart:]
	for _, row := range dataRows {
		price := models.StageInternationalMoveAccessorialPrice{
			Market:          getCell(row.Cells, firstColumnIndexStart),
			ServiceProvided: getCell(row.Cells, secondColumnIndexStart),
			PricePerUnit:    getCell(row.Cells, thirdColumnIndexStart),
		}

		// All the rows are consecutive, if we get a blank we're done
		if price.Market == "" {
			break
		}

		if params.ShowOutput == true {
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
	dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[additionalPricesRowIndexStart:]
	for _, row := range dataRows {
		price := models.StageDomesticInternationalAdditionalPrice{
			Market:       getCell(row.Cells, firstColumnIndexStart),
			ShipmentType: getCell(row.Cells, secondColumnIndexStart),
			Factor:       getCell(row.Cells, thirdColumnIndexStart),
		}

		// All the rows are consecutive, if we get a blank we're done
		if price.Market == "" {
			break
		}

		if params.ShowOutput == true {
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

	dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[domAccessorialRowIndexStart-2 : domAccessorialRowIndexStart-1]
	err := helperCheckHeadersFor5a("Services Schedule", "Service Provided", "PricePerUnitofMeasure", dataRows)
	if err != nil {
		return err
	}

	dataRows = params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[domAccessorialRowIndexStart-1 : domAccessorialRowIndexStart]
	err = helperCheckHeadersFor5a("X", "EXAMPLE (per unit of measure)", "$X.XX", dataRows)
	if err != nil {
		return err
	}

	dataRows = params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[intlAccessorialRowIndexStart-2 : intlAccessorialRowIndexStart-1]
	err = helperCheckHeadersFor5a("Market", "Service Provided", "PricePerUnitofMeasure", dataRows)
	if err != nil {
		return err
	}

	dataRows = params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[intlAccessorialRowIndexStart-1 : intlAccessorialRowIndexStart]
	err = helperCheckHeadersFor5a("X", "EXAMPLE (per unit of measure)", "$X.XX", dataRows)
	if err != nil {
		return err
	}

	dataRows = params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[additionalPricesRowIndexStart-2 : additionalPricesRowIndexStart-1]
	err = helperCheckHeadersFor5a("Market", "Shipment Type", "Factor", dataRows)
	if err != nil {
		return err
	}

	dataRows = params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[additionalPricesRowIndexStart-1 : additionalPricesRowIndexStart]
	return helperCheckHeadersFor5a("CONUS / OCONUS", "EXAMPLE", "X.XX", dataRows)
}

func helperCheckHeadersFor5a(firstHeader string, secondHeader string, thirdHeader string, dataRows []*xlsx.Row) error {
	const firstColumnIndexStart = 2
	const secondColumnIndexStart = 3
	const thirdColumnIndexStart = 4

	for _, dataRow := range dataRows {
		if header := getCell(dataRow.Cells, firstColumnIndexStart); header != firstHeader {
			return fmt.Errorf("verifyAccessAndAddPrices expected to find header '%s', but received header '%s'", firstHeader, header)
		}
		if header := getCell(dataRow.Cells, secondColumnIndexStart); header != secondHeader {
			return fmt.Errorf("verifyAccessAndAddPrices expected to find header '%s', but received header '%s'", secondHeader, header)
		}
		if header := removeWhiteSpace(getCell(dataRow.Cells, thirdColumnIndexStart)); header != thirdHeader {
			return fmt.Errorf("verifyAccessAndAddPrices expected to find header '%s', but received header '%s'", thirdHeader, header)
		}
	}
	return nil
}
