package pricing

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// parseServiceAreas: parser for: 1b) Service Areas
var parseDomesticServiceAreas processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 4         // 1b) Service Areas
	const serviceAreaRowIndexStart int = 9 // start at row 9 to get the service areas
	const basePointCityColumn int = 2
	const stateColumn int = 3
	const serviceAreaNumberColumn int = 4
	const zip3sColumn int = 5

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticServiceAreas expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageDomesticServiceArea")

	var domServAreas []models.StageDomesticServiceArea
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := serviceAreaRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		domServArea := models.StageDomesticServiceArea{
			BasePointCity:     mustGetCell(sheet, rowIndex, basePointCityColumn),
			State:             mustGetCell(sheet, rowIndex, stateColumn),
			ServiceAreaNumber: mustGetCell(sheet, rowIndex, serviceAreaNumberColumn),
			Zip3s:             mustGetCell(sheet, rowIndex, zip3sColumn),
		}
		// All the rows are consecutive, if we get to a blank one we're done
		if domServArea.BasePointCity == "" {
			break
		}

		prefixPrinter.Printf("%+v\n", domServArea)

		domServAreas = append(domServAreas, domServArea)
	}

	return domServAreas, nil
}

var parseInternationalServiceAreas processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 4         // 1b) Service Areas
	const serviceAreaRowIndexStart int = 9 // start at row 9 to get the service areas
	const internationalRateAreaColumn int = 9
	const rateAreaIDColumn int = 10

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseInternationalServiceAreas expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageInternationalServiceArea")

	var intlServAreas []models.StageInternationalServiceArea
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := serviceAreaRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		intlServArea := models.StageInternationalServiceArea{
			RateArea:   mustGetCell(sheet, rowIndex, internationalRateAreaColumn),
			RateAreaID: mustGetCell(sheet, rowIndex, rateAreaIDColumn),
		}
		// All the rows are consecutive, if we get to a blank one we're done
		if intlServArea.RateArea == "" {
			break
		}

		prefixPrinter.Printf("%+v\n", intlServArea)

		intlServAreas = append(intlServAreas, intlServArea)
	}

	return intlServAreas, nil
}

// verifyServiceAreas: verification for: 1b) Service Areas
var verifyServiceAreas verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 4         // 1b) Service Areas
	const serviceAreaRowIndexStart int = 9 // start at row 6 to get the headings
	const basePointCityColumn int = 2
	const stateColumn int = 3
	const serviceAreaNumberColumn int = 4
	const zip3sColumn int = 5
	const internationalRateAreaColumn int = 9
	const rateAreaIDColumn int = 10

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyServiceAreas expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	// Only check header of domestic and international service areas
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	dataRowsIndex := serviceAreaRowIndexStart - 1
	if header := mustGetCell(sheet, dataRowsIndex, basePointCityColumn); header != "Base Point City" {
		return fmt.Errorf("verifyServiceAreas expected to find header 'Base Point City', but received header '%s'", header)
	}
	if header := mustGetCell(sheet, dataRowsIndex, stateColumn); header != "State" {
		return fmt.Errorf("verifyServiceAreas expected to find header 'State', but received header '%s'", header)
	}
	if header := removeWhiteSpace(mustGetCell(sheet, dataRowsIndex, serviceAreaNumberColumn)); header != "ServiceAreaNumber" {
		return fmt.Errorf("verifyServiceAreas expected to find header 'ServiceAreaNumber', but received header '%s'", header)
	}
	if header := removeWhiteSpace(mustGetCell(sheet, dataRowsIndex, zip3sColumn)); header != "IncludedZip3's" {
		return fmt.Errorf("verifyServiceAreas expected to find header \"IncludedZip3's\", but received header '%s'", header)
	}
	if header := mustGetCell(sheet, dataRowsIndex, internationalRateAreaColumn); header != "International Rate Area" {
		return fmt.Errorf("verifyServiceAreas expected to find header 'International Rate Area', but received header '%s'", header)
	}
	if header := mustGetCell(sheet, dataRowsIndex, rateAreaIDColumn); header != "Rate Area ID" {
		return fmt.Errorf("verifyServiceAreas expected to find header 'Rate Area ID', but received header '%s'", header)
	}

	return nil
}
