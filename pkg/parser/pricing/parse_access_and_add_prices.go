package pricing

import (
	"fmt"

	"github.com/tealeg/xlsx/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

var parseDomesticMoveAccessorialPrices processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 17 // 5a) Access. and Add. Prices
	const domAccessorialRowIndexStart int = 11
	const firstColumnIndexStart = 2
	const secondColumnIndexStart = 3
	const thirdColumnIndexStart = 4

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticMoveAccessorialPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageDomesticMoveAccessorialPrice")

	var prices []models.StageDomesticMoveAccessorialPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := domAccessorialRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		price := models.StageDomesticMoveAccessorialPrice{
			ServicesSchedule: mustGetCell(sheet, rowIndex, firstColumnIndexStart),
			ServiceProvided:  mustGetCell(sheet, rowIndex, secondColumnIndexStart),
			PricePerUnit:     mustGetCell(sheet, rowIndex, thirdColumnIndexStart),
		}

		// All the rows are consecutive, if we get a blank we're done
		if price.ServicesSchedule == "" {
			break
		}

		prefixPrinter.Printf("%+v\n", price)

		prices = append(prices, price)
	}
	return prices, nil
}

var parseInternationalMoveAccessorialPrices processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 17 // 5a) Access. and Add. Prices
	const intlAccessorialRowIndexStart int = 25
	const firstColumnIndexStart = 2
	const secondColumnIndexStart = 3
	const thirdColumnIndexStart = 4

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseInternationalMoveAccessorialPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageInternationalMoveAccessorialPrice")

	var prices []models.StageInternationalMoveAccessorialPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := intlAccessorialRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		price := models.StageInternationalMoveAccessorialPrice{
			Market:          mustGetCell(sheet, rowIndex, firstColumnIndexStart),
			ServiceProvided: mustGetCell(sheet, rowIndex, secondColumnIndexStart),
			PricePerUnit:    mustGetCell(sheet, rowIndex, thirdColumnIndexStart),
		}

		// All the rows are consecutive, if we get a blank we're done
		if price.Market == "" {
			break
		}

		prefixPrinter.Printf("%+v\n", price)

		prices = append(prices, price)
	}
	return prices, nil
}

var parseDomesticInternationalAdditionalPrices processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 17 // 5a) Access. and Add. Prices
	const additionalPricesRowIndexStart int = 39
	const firstColumnIndexStart = 2
	const secondColumnIndexStart = 3
	const thirdColumnIndexStart = 4

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticInternationalAdditionalPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageDomesticInternationalAdditionalPrice")

	var prices []models.StageDomesticInternationalAdditionalPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := additionalPricesRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		price := models.StageDomesticInternationalAdditionalPrice{
			Market:       mustGetCell(sheet, rowIndex, firstColumnIndexStart),
			ShipmentType: mustGetCell(sheet, rowIndex, secondColumnIndexStart),
			Factor:       mustGetCell(sheet, rowIndex, thirdColumnIndexStart),
		}

		// All the rows are consecutive, if we get a blank we're done
		if price.Market == "" {
			break
		}

		prefixPrinter.Printf("%+v\n", price)

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

func helperCheckHeadersFor5a(firstHeader string, secondHeader string, thirdHeader string, sheet *xlsx.Sheet, dataRowsIndexBegin, dataRowsIndexEnd int) error {
	const firstColumnIndexStart = 2
	const secondColumnIndexStart = 3
	const thirdColumnIndexStart = 4

	for rowIndex := dataRowsIndexBegin; rowIndex < dataRowsIndexEnd; rowIndex++ {
		if header := mustGetCell(sheet, rowIndex, firstColumnIndexStart); header != firstHeader {
			return fmt.Errorf("verifyAccessAndAddPrices expected to find header '%s', but received header '%s'", firstHeader, header)
		}
		if header := mustGetCell(sheet, rowIndex, secondColumnIndexStart); header != secondHeader {
			return fmt.Errorf("verifyAccessAndAddPrices expected to find header '%s', but received header '%s'", secondHeader, header)
		}
		if header := removeWhiteSpace(mustGetCell(sheet, rowIndex, thirdColumnIndexStart)); header != thirdHeader {
			return fmt.Errorf("verifyAccessAndAddPrices expected to find header '%s', but received header '%s'", thirdHeader, header)
		}
	}
	return nil
}
